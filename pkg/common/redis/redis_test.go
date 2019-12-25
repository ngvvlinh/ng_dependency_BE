package redis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var store Store

func init() {
	redisAddress := "redis://redis:6379"
	store = NewWithPool(redisAddress)
}

func TestGetSetInterface(T *testing.T) {
	type Foo struct {
		Bar int
		Baz string
	}

	T.Run("Test set interface", func(t *testing.T) {
		err := store.Set("foo", &Foo{10, "sample"})
		require.Nil(t, err)
	})

	T.Run("Test get interface", func(t *testing.T) {
		var foo Foo
		err := store.Get("foo", &foo)
		require.Nil(t, err)
		require.Equal(t, 10, foo.Bar)
		require.Equal(t, "sample", foo.Baz)
	})

	T.Run("Test get interface (not found)", func(t *testing.T) {
		var foo Foo
		err := store.Get("not_found", &foo)
		require.EqualError(t, err, "redigo: nil returned")
	})

	T.Run("Test set interface with ttl", func(t *testing.T) {
		err := store.SetWithTTL("foo1", &Foo{Bar: 10, Baz: "sample"}, 2)
		require.Nil(t, err)
	})

	T.Run("Test get interface ttl", func(t *testing.T) {
		ttl, err := store.GetTTL("foo1")
		require.Nil(t, err)

		t.Log("ttl: ", ttl)
		require.True(t, (0 <= ttl) && (ttl <= 2) || (ttl == -2))
	})
}

func TestGetSetString(T *testing.T) {
	T.Run("Test set string", func(t *testing.T) {
		err := store.SetString("foo", "bar")
		require.Nil(t, err)
	})

	T.Run("Test get string", func(t *testing.T) {
		foo, err := store.GetString("foo")
		require.Nil(t, err)
		require.Equal(t, "bar", foo)
	})

	T.Run("Test string is existed", func(t *testing.T) {
		exist := store.IsExist("foo")
		require.True(t, exist)
	})

	T.Run("Test delete string", func(t *testing.T) {
		err := store.Del("foo")
		require.Nil(t, err)

		exist := store.IsExist("foo")
		require.False(t, exist)
	})

	T.Run("Test set string with ttl", func(t *testing.T) {
		err := store.SetStringWithTTL("string1", "string1", 2)
		require.Nil(t, err)
	})
}

func TestGetStrings(t *testing.T) {
	err := store.SetString("t:foo", "bar")
	require.Nil(t, err)
	err = store.SetString("t:bar", "baz")
	require.Nil(t, err)

	values, err := store.GetStrings("t:*")
	require.Nil(t, err)
	require.NotEmpty(t, values)
}

func TestGetSetUint64(T *testing.T) {
	T.Run("Test set uint64", func(t *testing.T) {
		err := store.SetUint64("ten", 10)
		require.Nil(t, err)
	})

	T.Run("Test get uint64", func(t *testing.T) {
		ten, err := store.GetUint64("ten")
		require.Nil(t, err)
		require.Equal(t, uint64(10), ten)
	})

	T.Run("Test set uint64 with ttl", func(t *testing.T) {
		err := store.SetUint64WithTTL("six", 6, 6)
		require.Nil(t, err)
	})
}

func TestDel(t *testing.T) {
	values, err := store.GetStrings("t:*")
	require.NoError(t, err)
	require.NotEmpty(t, values)

	err = store.Del(values...)
	require.NoError(t, err)

	values, err = store.GetStrings("t:*")
	require.NoError(t, err)
	require.Empty(t, values)
}
