package cmgrpc

import (
	"context"
	"fmt"
	"io"
	"net/http"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/twitchtv/twirp"
	google_rpc "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"etop.vn/backend/pkg/common/auth"
	"etop.vn/common/l"
)

var ll = l.New()

// Authorization definition
const (
	AuthorizationHeader = "authorization"
	AuthorizationType   = "bearer"
)

// AppendAccessToken ...
func AppendAccessToken(ctx context.Context, accessToken string) context.Context {
	headers := make(http.Header)
	headers.Set("Authorization", "Bearer "+accessToken)
	ctx, err := twirp.WithHTTPRequestHeaders(ctx, headers)
	if err != nil {
		ll.Panic("HTTP Headers", l.Error(err))
	}
	return ctx
}

// AccessTokenFromContext ...
func AccessTokenFromContext(ctx context.Context) string {
	token, err := grpc_auth.AuthFromMD(ctx, AuthorizationType)
	if err != nil {
		return ""
	}
	return token
}

// AppendMetadata ...
func AppendMetadata(ctx context.Context, pairs ...string) context.Context {
	return metadata.NewOutgoingContext(ctx, metadata.Pairs(pairs...))
}

// ForwardMetadata ...
func ForwardMetadata(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	return metadata.NewOutgoingContext(ctx, md)
}

// ForwardMetadataUnaryServerInterceptor ...
func ForwardMetadataUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		ctx = ForwardMetadata(ctx)
		return handler(ctx, req)
	}
}

// Status ...
func Status(code codes.Code, msg string) *google_rpc.Status {
	return &google_rpc.Status{
		Code:    int32(code),
		Message: msg,
	}
}

// Statusf ...
func Statusf(code codes.Code, format string, args ...interface{}) *google_rpc.Status {
	return &google_rpc.Status{
		Code:    int32(code),
		Message: fmt.Sprintf(format, args...),
	}
}

// AuthFunc ...
type AuthFunc func(ctx context.Context, fullMethod string) (context.Context, error)

// Authentication ...
func Authentication(authorizer auth.Authorizer, exceptions []string) AuthFunc {
	if authorizer == nil {
		ll.Panic("Nil validator")
	}
	return func(ctx context.Context, fullMethod string) (context.Context, error) {
		ignore := false
		for _, exception := range exceptions {
			if exception == fullMethod {
				ignore = true
				break
			}
		}
		tokenStr, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil && !ignore {
			ll.Warn("No authorization header", l.String("method", fullMethod), l.Error(err))
			return ctx, err
		}

		ctx, err = authorizer.Authorize(ctx, tokenStr)
		if err != nil && !ignore {
			ll.Warn("Invalid token", l.String("token", tokenStr), l.Error(err))
			return ctx, status.Errorf(codes.Unauthenticated, "Request unauthenticated")
		}
		return ctx, nil
	}
}

// AuthUnaryServerInterceptor ...
func AuthUnaryServerInterceptor(authFunc AuthFunc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := authFunc(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(newCtx, req)
	}
}

// IDGenerator ...
type IDGenerator func() string

// Append new correlation-id for outgoing context
func generateCorrelationID(ctx context.Context, reqIDGen IDGenerator) context.Context {
	const header = "correlation-id"
	var reqID string
	inMD, _ := metadata.FromIncomingContext(ctx)
	if ids, ok := inMD[header]; ok && len(ids) > 0 {
		reqID = ids[0]
	} else {
		reqID = reqIDGen()
	}
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(header, reqID))
	return ctx
}

// RequestLogger ...
type RequestLogger func(ctx context.Context, method string, req, resp interface{}, err error)

// CensorRequest censors sensitive data
func CensorRequest(req interface{}) {
	if req, ok := req.(interface {
		Censor()
	}); ok {
		req.Censor()
	}
}

// ResponseHasError ...
func ResponseHasError(resp interface{}) (isError bool) {
	if resp, ok := resp.(interface {
		HasError() bool
	}); ok {
		return resp.HasError()
	}
	return false
}

var jSON = runtime.JSONPb{
	OrigName:     true,
	EmitDefaults: true,
}

// NewJSONDecoder ...
func NewJSONDecoder(r io.Reader) runtime.Decoder {
	return jSON.NewDecoder(r)
}

// NewJSONEncoder ...
func NewJSONEncoder(w io.Writer) runtime.Encoder {
	return jSON.NewEncoder(w)
}

// MarshalJSON encodes JSON in compatible with GRPC
func MarshalJSON(v interface{}) ([]byte, error) {
	return jSON.Marshal(v)
}

// UnmarshalJSON decodes JSON in compatible with GRPC
func UnmarshalJSON(data []byte, v interface{}) error {
	return jSON.Unmarshal(data, v)
}

// OptToken ...
func OptToken(accessToken string) grpc.CallOption {
	return grpc.Header(&metadata.MD{
		AuthorizationHeader: []string{AuthorizationType + " " + accessToken},
	})
}
