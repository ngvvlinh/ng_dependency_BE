package tokens

import (
	"context"
	"time"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/authorization/auth"
	"o.o/backend/pkg/common/redis"
	"o.o/backend/pkg/etop/authorize/claims"
	"o.o/common/l"
)

var ll = l.New()

const (
	UsageAccessToken = "AT"

	// DefaultAccessTokenTTL ...
	DefaultAccessTokenTTL = 60 * 60 * 24 * 7 // 3 days
)

type Validator interface {
	Validate(tokenStr string) (*claims.Claim, error)
}

// TokenStore ...
type TokenStore struct {
	auth auth.Generator
}

func NewTokenStore(r redis.Store) TokenStore {
	return TokenStore{auth: auth.NewGenerator(r)}
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

func (s TokenStore) GenerateToken(ctx context.Context, cmd *GenerateTokenCommand) error {
	claim := &claims.Claim{
		ClaimInfo: cmd.ClaimInfo,
	}
	tok, err := s.generateWithClaim(claim, cmd.TTL)
	if err != nil {
		return err
	}
	cmd.Result = tok
	return nil
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
	v.CAS = cm.RandomInt64()
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
	v.CAS = cm.RandomInt64()
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
