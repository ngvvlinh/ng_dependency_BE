package auth

import (
	"os"
	"testing"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/common/redis"
)

var gFoo Generator

func TestMain(M *testing.M) {
	endpoint := "redis:6379"
	redisPool := &redigo.Pool{
		MaxIdle:     50,
		MaxActive:   0,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", endpoint)
			if err != nil {
				ll.Fatal("Redis: unable to dial on %s", l.Error(err), l.String("endpoint", endpoint))
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			ll.Fatal("err", l.Error(err))
			return err
		},
	}
	rStore := redis.New(redisPool)
	gFoo = NewGenerator(rStore)

	os.Exit(M.Run())
}

func TestGenerateValidateRevokeToken(t *testing.T) {
	id := int64(123456789)
	var tokenStr string

	t.Run("Generate", func(t *testing.T) {
		tok, err := gFoo.Generate(UsageAccessToken, id, DefaultAccessTokenTTL)
		require.NoError(t, err)

		tokenStr = tok.TokenStr
	})

	t.Run("Validate", func(t *testing.T) {
		tok, err := gFoo.Validate(UsageAccessToken, tokenStr, nil)
		require.NoError(t, err)

		assert.Equal(t, id, tok.UserID)
	})

	t.Run("Revoke", func(t *testing.T) {
		err := gFoo.Revoke(UsageAccessToken, tokenStr)
		require.NoError(t, err)

		_, err = gFoo.Validate(UsageAccessToken, tokenStr, nil)
		assert.EqualError(t, err, "invalid token")
	})
}

func TestGenerateWithValueToken(T *testing.T) {
	id := int64(123456789)
	t, err := gFoo.GenerateWithValue(&Token{
		Usage:  UsageAccessToken,
		UserID: id,
		Value: map[string]string{
			"hello": "foo",
		},
	}, DefaultAccessTokenTTL)
	if err != nil {
		T.Fatal(err)
	}

	defer func() {
		if err := gFoo.Revoke(UsageAccessToken, t.TokenStr); err != nil {
			panic(err)
		}
	}()

	var v map[string]string
	_, err = gFoo.Validate(UsageAccessToken, t.TokenStr, &v)
	if err != nil {
		T.Fatal(err)
	}

	if v["hello"] != "foo" {
		T.Fatal("Token value was not set correctly")
	}
}

func TestTokenWithTTL(T *testing.T) {
	id := int64(123456789)
	t, err := gFoo.Generate(UsageAccessToken, id, 1)
	if err != nil {
		T.Fatal(err)
	}

	// Make sure TTL passed
	time.Sleep(1005 * time.Millisecond)

	if _, err := gFoo.Validate(UsageAccessToken, t.TokenStr, nil); err == nil {
		defer func() {
			if err := gFoo.Revoke(UsageAccessToken, t.TokenStr); err != nil {
				panic(err)
			}
		}()

		T.Fatal("Key does not expired")
	}
}
