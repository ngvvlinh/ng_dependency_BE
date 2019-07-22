package tokens

import (
	"context"
	"time"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/auth"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/etop/authorize/claims"
	"etop.vn/common/bus"
	"etop.vn/common/l"
)

var ll = l.New()

const (
	UsageAccessToken = "AT"

	// DefaultAccessTokenTTL ...
	DefaultAccessTokenTTL = 60 * 60 * 24 * 7 // 3 days
)

var Store TokenStore

func init() {
	bus.AddHandlers("session",
		GenerateToken,
		UpdateSession,
	)
}

// InitTokenStore ...
func Init(r redis.Store) {
	Store = TokenStore{
		auth: auth.NewGenerator(r),
	}
}

// TokenStore ...
type TokenStore struct {
	auth auth.Generator
}

type GenerateTokenCommand struct {
	claims.ClaimInfo
	TTL int

	Result *auth.Token
}

type UpdateSessionCommand struct {
	Token  string
	Values map[string]string
}

func GenerateToken(ctx context.Context, cmd *GenerateTokenCommand) error {
	claim := &claims.Claim{
		ClaimInfo: cmd.ClaimInfo,
	}
	tok, err := Store.generateWithClaim(claim, cmd.TTL)
	if err != nil {
		return err
	}
	cmd.Result = tok
	return nil
}

func UpdateSession(ctx context.Context, cmd *UpdateSessionCommand) error {
	return Store.UpdateSession(cmd.Token, cmd.Values)
}

// GenerateWithClaim ...
func (s TokenStore) generateWithClaim(v *claims.Claim, ttl int) (*auth.Token, error) {
	t := &auth.Token{
		Usage:  UsageAccessToken,
		UserID: v.UserID,
		Value:  v,
	}
	if ttl == 0 {
		ttl = DefaultAccessTokenTTL
	}
	v.CAS = cm.NewID()
	t, err := s.auth.GenerateWithValue(t, ttl)
	return t, err
}

// Validate ...
func (s TokenStore) Validate(tokenStr string) (*claims.Claim, error) {
	var v claims.Claim
	tok, err := s.auth.Validate(UsageAccessToken, tokenStr, &v)
	if err != nil {
		ll.Error("Invalid access token", l.Error(err))
		return nil, err
	}

	// UpdateInfo LastLoginAt every 1 minute
	now := time.Now()
	if now.Sub(v.LastLoginAt) > 1*time.Minute {
		vv := v
		vv.LastLoginAt = now
		tok.Value = vv
		err = s.auth.SetTTL(tok, DefaultAccessTokenTTL)
		if err != nil {
			ll.Error("Unable to update TTL", l.Error(err))
		}
	}

	v.Token = tokenStr
	v.UserID = tok.UserID
	return &v, nil
}

func (s TokenStore) UpdateSession(tokStr string, values map[string]string) error {
	if tokStr == "" {
		return cm.Errorf(cm.Internal, nil, "no token")
	}

	var v claims.Claim
	current, err := s.auth.Validate(UsageAccessToken, tokStr, &v)
	if err != nil {
		ll.Error("Invalid access token", l.Error(err))
		return err
	}
	v.CAS = cm.NewID()
	if v.Extra == nil {
		v.Extra = make(map[string]string)
	}
	for key, val := range values {
		v.Extra[key] = val
	}
	current.Value = v

	// TODO: refactor (should not expose SetInfo and current.ToValue)
	return s.auth.SetInfo(current.ToKey(), current.ToValue())
}
