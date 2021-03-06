package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"

	cm "o.o/backend/pkg/common"
	cmWrapper "o.o/backend/pkg/common/apifw/wrapper"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/headers"
	"o.o/backend/pkg/etop/authorize/permission"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/httprpc"
	"o.o/common/jsonx"
	"o.o/common/l"
	"o.o/common/xerrors"
)

var (
	ll      = l.New()
	decoder = schema.NewDecoder()
)

func init() {
	decoder.SetAliasTag("json")
}

type Handler func(c *Context) error
type MiddlewareFunc func(Handler) Handler

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter
	session.Session
	httprouter.Params

	hasResult bool
	rawResult bool
	result    interface{}
}

func (c *Context) Context() context.Context {
	return c.Req.Context()
}

func (c *Context) SetResultRaw() http.ResponseWriter {
	if c.hasResult {
		ll.Panic("Must only set result once!")
	}
	c.rawResult = true
	return c.Resp
}

func (c *Context) SetResult(v interface{}) {
	if c.hasResult {
		ll.Panic("Must only set result once!")
	}
	if v == nil {
		ll.Panic("Result is empty")
	}
	vv := reflect.ValueOf(v)
	switch {
	case vv.Kind() == reflect.Ptr && !vv.IsNil():
	case vv.Kind() == reflect.Map:
	default:
		ll.Panic("Result must be map or pointer to struct")
	}
	c.result = v
	c.hasResult = true
}

func (c *Context) DecodeJson(v interface{}) error {
	body, err := ioutil.ReadAll(c.Req.Body)
	if err != nil {
		return cm.Error(cm.InvalidArgument, err.Error(), err)
	}
	err = jsonx.Unmarshal(body, v)
	if err != nil {
		return cm.Error(cm.InvalidArgument, err.Error(), err)
	}

	ll.Info("->"+c.Req.URL.Path, l.String("data", string(body)))
	return nil
}

func (c *Context) GetBody(v interface{}) ([]byte, error) {
	defer c.Req.Body.Close()
	body, err := ioutil.ReadAll(c.Req.Body)
	if err != nil {
		return nil, cm.Error(cm.InvalidArgument, err.Error(), err)
	}
	return body, nil
}

const defaultMemory = 32 << 20 // 32 MB

func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Req.ParseMultipartForm(defaultMemory)
	return c.Req.MultipartForm, err
}

func (c *Context) DecodeFormUrlEncoded(v interface{}) error {
	if err := c.Req.ParseForm(); err != nil {
		return cm.Error(cm.InvalidArgument, err.Error(), err)
	}
	ll.Info("->"+c.Req.URL.Path, l.String("data", c.Req.Form.Encode()))

	err := decoder.Decode(v, c.Req.Form)
	if err != nil {
		return cm.Error(cm.InvalidArgument, err.Error(), err)
	}
	return nil
}

type Options struct {
	LogRequest bool
}

type Router struct {
	*httprouter.Router

	middlewares []func(Handler) Handler
}

func New() *Router {
	r := &httprouter.Router{}
	return &Router{Router: r}
}

func (rt *Router) Use(m func(Handler) Handler) {
	rt.middlewares = append(rt.middlewares, m)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (rt *Router) GET(path string, handle Handler, ms ...MiddlewareFunc) {
	rt.Handle("GET", path, handle, ms...)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (rt *Router) HEAD(path string, handle Handler, ms ...MiddlewareFunc) {
	rt.Handle("HEAD", path, handle, ms...)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (rt *Router) OPTIONS(path string, handle Handler, ms ...MiddlewareFunc) {
	rt.Handle("OPTIONS", path, handle, ms...)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (rt *Router) POST(path string, handle Handler, ms ...MiddlewareFunc) {
	rt.Handle("POST", path, handle, ms...)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (rt *Router) PUT(path string, handle Handler, ms ...MiddlewareFunc) {
	rt.Handle("PUT", path, handle, ms...)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (rt *Router) PATCH(path string, handle Handler, ms ...MiddlewareFunc) {
	rt.Handle("PATCH", path, handle, ms...)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (rt *Router) DELETE(path string, handle Handler, ms ...MiddlewareFunc) {
	rt.Handle("DELETE", path, handle, ms...)
}

func (rt *Router) RawHandle(method, path string, h http.Handler) {
	rt.Router.Handle(method, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h.ServeHTTP(w, r)
	})
}

func (rt *Router) Handle(method, path string, h Handler, ms ...MiddlewareFunc) {
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i](h)
	}
	for i := len(rt.middlewares) - 1; i >= 0; i-- {
		m := rt.middlewares[i]
		h = m(h)
	}
	rt.Router.Handle(method, path, rt.wrapJSON(h))
}

func (rt *Router) wrapJSON(next Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		c := &Context{
			Req:    r,
			Resp:   w,
			Params: params,
		}
		var err error
		defer func() {
			// The inner logic is responsible to write the response
			if c.rawResult {
				if err != nil {
					ll.Warn("error with custom content-type", l.Error(err), l.Any("header", w.Header()))
				}
				return
			}

			if err != nil {
				twerr := xerrors.TwirpError(err)
				statusCode := httprpc.ServerHTTPStatusFromErrorCode(twerr.Code())
				jerr := xerrors.ToErrorJSON(twerr)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(statusCode)
				_ = json.NewEncoder(w).Encode(jerr)
				return
			}

			switch {
			case c.result != nil:
				var respBytes []byte
				if respBytes, err = jsonx.Marshal(c.result); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(&xerrors.ErrorJSON{
						Code: cm.Internal.String(),
						Msg:  "failed to marshal json response",
					})
					ll.Error("Failed to marshal json response", l.Error(err))
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(respBytes)

			default:
				ll.Panic("If error is nil, a result must always be provided!")
			}
		}()

		err = next(c)
	}
}

func RecoverAndLog(logRequest bool) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(c *Context) (_err error) {
			t0 := time.Now()
			ctx := bus.NewRootContext(c.Req.Context())
			req := c.Req.WithContext(ctx)
			c.Req = req

			var reqData []byte
			defer func() {
				d := time.Since(t0)
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					_err = cm.Errorf(cm.RuntimePanic, err, "")
					ll.Error("->"+req.RequestURI+" (recovered)",
						l.Int("s", 500),
						l.Duration("d", d),
						l.String("req", string(reqData)),
						l.Error(_err))
					fmt.Printf("%+v", _err)
					return
				}
				if parseErr, ok := _err.(*xerrors.APIError); ok {
					result := parseErr.Meta["result"]
					if result == "ignore" {
						twError := xerrors.TwirpError(_err)
						go cmWrapper.SendErrorToBot(ctx, req.RequestURI, cmWrapper.AdaptSession(c.Session), reqData, twError, nil, d, xerrors.LevelPartialError, nil)
						c.result = map[string]string{
							"code": "ok",
						}
						_err = nil
					}
				}
				if _err == nil {
					if errs := cmWrapper.HasErrors(c.result); errs != nil {
						ll.Warn("->"+req.RequestURI,
							l.Duration("d", d),
							l.String("req", string(reqData)),
							l.String("resp", jsonx.MustMarshalToString(c.result)))
						go cmWrapper.SendErrorToBot(ctx, req.RequestURI, cmWrapper.AdaptSession(c.Session), reqData, nil, errs, d, xerrors.LevelPartialError, nil)
						return
					}

					ll.Debug("->"+req.RequestURI,
						l.Int("s", http.StatusOK),
						l.Duration("d", d),
						l.String("req", string(reqData)))
					return
				}

				lvl := xerrors.GetTraceLevel(_err)
				if lvl <= xerrors.LevelTrival {
					if cmenv.IsDev() {
						ll.Warn("->"+req.RequestURI,
							l.Duration("d", d),
							l.String("req", string(reqData)),
							l.Error(_err))
					}
					return
				}

				twError := xerrors.TwirpError(_err)
				ll.Error("->"+req.RequestURI,
					l.Duration("d", d),
					l.String("req", string(reqData)),
					l.Error(_err))
				if lvl >= xerrors.LevelTrace {
					cmWrapper.PrintErrorWithStack(ctx, _err, nil)
				}
				go cmWrapper.SendErrorToBot(ctx, req.RequestURI, cmWrapper.AdaptSession(c.Session), reqData, twError, nil, d, lvl, nil)
			}()

			if logRequest {
				var err error
				reqData, err = ioutil.ReadAll(req.Body)
				if err != nil {
					return err
				}
				req.Body = ioutil.NopCloser(bytes.NewReader(reqData))
			} else {
				reqData = []byte(`<not logged>`)
			}
			return next(c)
		}
	}
}

func Auth(perm permission.Decl, _ss session.Session) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(c *Context) error {
			c.Session = _ss // clone the session
			ctx := c.Context()
			tokenStr := headers.GetBearerTokenFromCtx(ctx)
			ctx, err := c.StartSession(ctx, perm, tokenStr)
			if err != nil {
				return err
			}
			return next(c)
		}
	}
}
