package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strconv"
	"strings"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

const (
	// DefaultTokenLength ...
	DefaultTokenLength = 32 // In bytes

	// DefaultAccessTokenTTL ttl in seconds
	DefaultAccessTokenTTL = 60 * 60 * 24 * 30

	UsageAccessToken       = "AT"
	UsageResetPassword     = "RP"
	UsageSToken            = "ST"
	UsageEmailVerification = "EV"
	UsagePhoneVerification = "PV"
	UsageRequestLogin      = "RL"
	UsageRegister          = "RT"

	UsagePartnerIntegration = "PI"
	UsageInviteUser         = "iv"
)

var ll = l.New()

// Validator interface
type Validator interface {
	Validate(usage, tokenStr string, v interface{}) (*Token, error)
}

// ApiKey represents a token used in request/response
type Token struct {
	TokenStr  string
	Usage     string
	UserID    dot.ID
	Value     interface{}
	ExpiresIn int
}

// Store interface contains methods
// to perform token actions
type Store interface {
	Generate(usage string, userID dot.ID, ttl int) (*Token, error)
	GenerateWithValue(tok *Token, ttl int) (*Token, error)
	Revoke(usage, tokenStr string) error
	SetTTL(token *Token, ttl int) error
	SetInfo(tokenStr string, value string) error
}

// Generator interface
type Generator interface {
	Validator
	Store
}

type generator struct {
	redisStore redis.Store
}

// NewGenerator returns new token generator
func NewGenerator(r redis.Store) Generator {
	return &generator{
		redisStore: r,
	}
}

// ToKey returns string that can be used as redis key
// from given token
func (t *Token) ToKey() string {
	return t.Usage + ":" + t.TokenStr
}

// ToValue returns string that can be used as redis value
// from given token
func (t *Token) ToValue() string {
	if t.Value == nil {
		return strconv.FormatInt(int64(t.UserID), 16)
	}

	data, err := jsonx.Marshal(t.Value)
	if err != nil {
		ll.Panic("Unable to marshal json", l.Error(err))
	}
	return strconv.FormatInt(int64(t.UserID), 16) + ":" + string(data)
}

// Generate creates token for given userID and TTL.
func (g *generator) Generate(usage string, userID dot.ID, ttl int) (*Token, error) {
	t := &Token{
		Usage:  usage,
		UserID: userID,
	}
	return g.generate(t, ttl)
}

// GenerateWithValue creates token with given value for given userID and TTL.
func (g *generator) GenerateWithValue(t *Token, ttl int) (*Token, error) {
	t.ExpiresIn = ttl
	return g.generate(t, ttl)
}

func (g *generator) generate(t *Token, ttl int) (*Token, error) {
	value := t.ToValue()

	if t.TokenStr != "" {
		key := t.ToKey()
		if g.redisStore.IsExist(key) {
			return t, errors.New("key already exists")
		}
		err := g.redisStore.SetStringWithTTL(key, value, ttl)
		return t, err
	}

	for retry := 0; retry < 3; retry++ {
		t.TokenStr = RandomToken(DefaultTokenLength)
		key := t.ToKey()
		if !g.redisStore.IsExist(key) {
			err := g.redisStore.SetStringWithTTL(key, value, ttl)
			return t, err
		}
	}

	panic("Unable to generate token, retried 3 times!")
}

func (g *generator) Validate(usage, token string, v interface{}) (*Token, error) {
	t := &Token{
		Usage:    usage,
		TokenStr: token,
	}

	// Check if the token exist in database
	key := t.ToKey()
	storedValue, err := g.redisStore.GetString(key)
	if err == redis.ErrNil || storedValue == "" {
		return t, cm.Errorf(cm.NotFound, nil, "invalid token")
	}
	if err != nil {
		return t, err
	}

	s := strings.SplitN(storedValue, ":", 2)
	userID, err := strconv.ParseInt(s[0], 16, 64)
	if err != nil {
		return t, cm.Errorf(cm.NotFound, err, "invalid token")
	}

	t.UserID = dot.ID(userID)
	if v != nil && len(s) > 1 && len(s[1]) > 0 {
		err = jsonx.Unmarshal([]byte(s[1]), v)
		if err != nil {
			ll.Error("Unable to decode json", l.Error(err), l.String("json", s[1]))
			return t, cm.Errorf(cm.NotFound, err, "invalid token")
		}
	}

	t.Value = v
	return t, nil
}

// Revoke deletes token from redis store.
func (g *generator) Revoke(usage, tokenStr string) error {
	t := Token{
		TokenStr: tokenStr,
		Usage:    usage,
	}
	key := t.ToKey()
	err := g.redisStore.Del(key)
	if err != nil {
		ll.Error("Error revoking token", l.Error(err))
	}
	return err
}

// SetInfo set information for given token
func (g *generator) SetInfo(tokenStr string, value string) error {
	ttl, err := g.redisStore.GetTTL(tokenStr)
	if err != nil {
		return err
	}

	err = g.redisStore.SetStringWithTTL(tokenStr, value, ttl)
	if err != nil {
		return err
	}

	return nil
}

// GetInfo return infomation for given token
func (g *generator) GetInfo(tokenStr string) (string, error) {
	storedValue, err := g.redisStore.GetString(tokenStr)
	if err != nil {
		return "", err
	}

	return storedValue, nil
}

// SetTTL ...
func (g *generator) SetTTL(t *Token, ttl int) error {
	return g.redisStore.SetStringWithTTL(t.ToKey(), t.ToValue(), ttl)
}

// RandomToken generate new base64 string from random byte array with given length
func RandomToken(length int) string {
	token := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, token); err != nil {
		panic(err)
	}
	return base64.RawURLEncoding.EncodeToString(token)
}
