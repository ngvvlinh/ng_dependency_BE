package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// Constants
const (
	AuthorizationHeader = "Authorization"
	AuthorizationScheme = "Bearer"
)

// Authorizer interface
type Authorizer interface {
	Authorize(ctx context.Context, tokenStr string) (context.Context, error)
}

// MultiAuthorizer ...
type MultiAuthorizer []Authorizer

// FromHTTPHeader ...
func FromHTTPHeader(h http.Header) (string, error) {
	s := h.Get(AuthorizationHeader)
	if s == "" {
		return "", nil
	}
	return FromHeaderString(s)
}

// FromHeaderString ...
func FromHeaderString(s string) (string, error) {
	splits := strings.SplitN(s, " ", 2)
	if len(splits) < 2 {
		return "", errors.New("bad authorization string")
	}
	if splits[0] != "Bearer" && splits[0] != "bearer" {
		return "", errors.New("request unauthenticated with " + AuthorizationScheme)
	}
	return splits[1], nil
}

// Authorize ...
func (authorizers MultiAuthorizer) Authorize(ctx context.Context, tokenStr string) (context.Context, error) {
	var (
		err    error
		newCtx context.Context
	)
	for _, f := range authorizers {
		newCtx, err = f.Authorize(ctx, tokenStr)
		if err != nil {
			continue
		}
		return newCtx, nil
	}
	return newCtx, err
}

// AuthorizerFunc implements Authorizer interface
type AuthorizerFunc func(context.Context, string) (context.Context, error)

// Authorize implements Authorizer interface
func (f AuthorizerFunc) Authorize(ctx context.Context, tokenStr string) (context.Context, error) {
	return f(ctx, tokenStr)
}
