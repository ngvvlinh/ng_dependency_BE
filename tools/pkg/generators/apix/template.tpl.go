package apix

const tplText = `
{{range $s := .Services}}

type {{.Name}}ServiceServer struct {
	{{.Name}}API
}

func New{{.Name}}ServiceServer(svc {{.Name}}API) Server{
	return &{{.Name}}ServiceServer {
		{{.Name}}API: svc,
	}
}

const {{.Name}}ServicePathPrefix = "/api{{.APIPath}}/"

func (s *{{$s.Name}}ServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	if req.Method != "POST" {
		return
	}
	switch req.URL.Path {
	{{range $m := .Methods}}
	case "/api{{$s.APIPath}}/{{.Name}}":
		s.serve{{.Name}}(ctx,resp,req)
		return
	{{end}}
	default:
		return
	}
}
{{range $m := .Methods}}
func (s *{{$s.Name}}ServiceServer) serve{{.Name}}(ctx context.Context, resp http.ResponseWriter, req *http.Request)  {
	header := req.Header.Get("Content-Type")
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	switch strings.TrimSpace(strings.ToLower(header[:i])) {
	case "application/json":
		s.serve{{.Name}}JSON(ctx, resp, req)
	default:
		msg := fmt.Sprintf("unexpected Content-Type: %q", req.Header.Get("Content-Type"))
		twerr := badRouteError(msg, req.Method, req.URL.Path)
		writeError(ctx, resp, twerr)
		return
	}
}

func (s *{{$s.Name}}ServiceServer) serve{{.Name}}JSON(ctx context.Context, resp http.ResponseWriter, req *http.Request) {
	var err error
	reqContent := new({{.Request.Type|type}})
	unmarshaler := jsonpb.Unmarshaler{AllowUnknownFields: true}
	if err := unmarshaler.Unmarshal(req.Body, reqContent); err != nil {
		writeError(ctx, resp, malformedRequestError("the json request could not be decoded").WithMeta("cause", err.Error()))
		return
	}
	// Call service method
	var respContent *{{.Response.Type|type}}
	func () {
		defer ensurePanicResponses(ctx, resp)
		respContent, err = s.{{$s.Name}}API.{{.Name}}(ctx, reqContent)
	}()
	if err != nil {
		writeError(ctx, resp, err)
		return
	}
	if respContent == nil {
		writeError(ctx, resp, twirp.InternalError("received a nil response"))
		return
	}
	var buf bytes.Buffer
	marshaler := &jsonpb.Marshaler{OrigName: true, EmitDefaults: true}
	if err = marshaler.Marshal(&buf, respContent); err != nil {
		writeError(ctx, resp, wrapInternal(err, "failed to marshal json response"))
		return
	}
	respBytes := buf.Bytes()
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBytes)))
	resp.WriteHeader(http.StatusOK)
	if  _,err := resp.Write(respBytes); err != nil {
		return
	}
}
{{end}}

func (s *{{$s.Name}}ServiceServer)PathPrefix() string {
	return {{.Name}}ServicePathPrefix
}
{{end}}
type Server interface {
	http.Handler
	PathPrefix() string
}

// ensurePanicResponses makes sure that rpc methods causing a panic still result in a Twirp Internal
// error response (status 500), and error hooks are properly called with the panic wrapped as an error.
// The panic is re-raised so it can be handled normally with middleware.
func ensurePanicResponses(ctx context.Context, resp http.ResponseWriter) {
	if r := recover(); r != nil {
		// Wrap the panic as an error so it can be passed to error hooks.
		// The original error is accessible from error hooks, but not visible in the response.
		err := errFromPanic(r)
		twerr := &internalWithCause{msg: "Internal service panic", cause: err}
		// Actually write the error
		writeError(ctx, resp, twerr)
		// If possible, flush the error to the wire.
		f, ok := resp.(http.Flusher)
		if ok {
			f.Flush()
		}
		panic(r)
	}
}

// errFromPanic returns the typed error if the recovered panic is an error, otherwise formats as error.
func errFromPanic(p interface{}) error {
	if err, ok := p.(error); ok {
		return err
	}
	return fmt.Errorf("panic: %v", p)
}

// internalWithCause is a Twirp Internal error wrapping an original error cause, accessible
// by github.com/pkg/errors.Cause, but the original error message is not exposed on Msg().
type internalWithCause struct {
	msg   string
	cause error
}

func (e *internalWithCause) Cause() error                                { return e.cause }
func (e *internalWithCause) Error() string                               { return e.msg + ": " + e.cause.Error() }
func (e *internalWithCause) Code() twirp.ErrorCode                       { return twirp.Internal }
func (e *internalWithCause) Msg() string                                 { return e.msg }
func (e *internalWithCause) Meta(key string) string                      { return "" }
func (e *internalWithCause) MetaMap() map[string]string                  { return nil }
func (e *internalWithCause) WithMeta(key string, val string) twirp.Error { return e }

// malformedRequestError is used when the twirp server cannot unmarshal a request
func malformedRequestError(msg string) twirp.Error {
	return twirp.NewError(twirp.Malformed, msg)
}

// badRouteError is used when the twirp server cannot route a request
func badRouteError(msg string, method, url string) twirp.Error {
	err := twirp.NewError(twirp.BadRoute, msg)
	err = err.WithMeta("twirp_invalid_route", method+" "+url)
	return err
}

// writeError writes Twirp errors in the response and triggers hooks.
func writeError(ctx context.Context, resp http.ResponseWriter, err error) {
		// Non-twirp errors are wrapped as Internal (default)
		twerr, ok := err.(twirp.Error)
		if !ok {
			twerr = twirp.InternalErrorWith(err)
		}

		statusCode := twirp.ServerHTTPStatusFromErrorCode(twerr.Code())

		respBody := marshalErrorToJSON(twerr)
		resp.Header().Set("Content-Type", "application/json") // Error responses are always JSON
		resp.Header().Set("Content-Length", strconv.Itoa(len(respBody)))
		resp.WriteHeader(statusCode) // set HTTP status code and send response

		_, writeErr := resp.Write(respBody)
		if writeErr != nil {
		// We have three options here. We could log the error, call the Error
		// hook, or just silently ignore the error.
		//
		// Logging is unacceptable because we don't have a user-controlled
		// logger; writing out to stderr without permission is too rude.
		//
		// Calling the Error hook would confuse users: it would mean the Error
		// hook got called twice for one request, which is likely to lead to
		// duplicated log messages and metrics, no matter how well we document
		// the behavior.
		//
		// Silently ignoring the error is our least-bad option. It's highly
		// likely that the connection is broken and the original 'err' says
		// so anyway.
		_ = writeErr
	}
}

// wrapInternal wraps an error with a prefix as an Internal error.
// The original error cause is accessible by github.com/pkg/errors.Cause.
func wrapInternal(err error, prefix string) twirp.Error {
	return twirp.InternalErrorWith(&wrappedError{prefix: prefix, cause: err})
}

type wrappedError struct {
	prefix string
	cause  error
}

func (e *wrappedError) Cause() error  { return e.cause }
func (e *wrappedError) Error() string { return e.prefix + ": " + e.cause.Error() }

// JSON serialization for errors
type twerrJSON struct {
	Code string            "json:\"code\""
	Msg  string            "json:\"msg\""
	Meta map[string]string "json:\"meta,omitempty\""
}

// marshalErrorToJSON returns JSON from a twirp.Error, that can be used as HTTP error response body.
// If serialization fails, it will use a descriptive Internal error instead.
func marshalErrorToJSON(twerr twirp.Error) []byte {
	// make sure that msg is not too large
	msg := twerr.Msg()
	if len(msg) > 1e6 {
		msg = msg[:1e6]
	}

	tj := twerrJSON{
		Code: string(twerr.Code()),
		Msg:  msg,
		Meta: twerr.MetaMap(),
	}

	buf, err := json.Marshal(&tj)
	if err != nil {
		buf = []byte("{\"type\": \"" + twirp.Internal + "\", \"msg\": \"There was an error but it could not be serialized into JSON\"}") // fallback
	}
	return buf
}
`
