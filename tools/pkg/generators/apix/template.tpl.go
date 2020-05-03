package apix

const tplText = `
func init() {
	httprpc.Register(NewServer)
}

func NewServer(builder interface{}, hooks ...*httprpc.Hooks) (httprpc.Server, bool) {
	switch builder := builder.(type) {
	{{range $s := .Services -}}
	case func() {{.Name}}Service:
		return New{{.Name}}ServiceServer(builder, hooks...), true
	{{end -}}
	default:
	{{if .Services|len|eq 0 -}}
		_ = builder
	{{end -}}
		return nil, false
	}
}

{{range $s := .Services}}
type {{.Name}}ServiceServer struct {
	hooks   httprpc.Hooks
	builder func() {{.Name}}Service
}

func New{{.Name}}ServiceServer(builder func() {{.Name}}Service, hooks ... *httprpc.Hooks) httprpc.Server{
	return &{{.Name}}ServiceServer {
		hooks: httprpc.WrapHooks(httprpc.ChainHooks(hooks...)),
		builder: builder,
	}
}

const {{.Name}}ServicePathPrefix = "{{.APIPath}}/"

func (s *{{$s.Name}}ServiceServer) PathPrefix() string {
	return {{.Name}}ServicePathPrefix
}

func (s *{{$s.Name}}ServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx, info := req.Context(), httprpc.HookInfo{Route: req.URL.Path, HTTPRequest: req}
	ctx, err := s.hooks.BeforeRequest(ctx, info)
	if err != nil {
		httprpc.WriteError(ctx, resp, s.hooks, info, err)
		return
	}
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, s.hooks, info, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, s.hooks, info, err)
		return
	}
	serve(ctx, resp, req, s.hooks, info, reqMsg, exec)
}

func (s *{{$s.Name}}ServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
{{range $m := .Methods -}}
	case "{{$s.APIPath}}/{{.Name}}":
	msg := {{(index .Request.Items 0).Type|new}}
	fn := func(ctx context.Context) (capi.Message, error) {
		inner := s.builder()
		ctx, err := s.hooks.BeforeServing(ctx, httprpc.HookInfo{Route: path, Request: msg}, inner)
		if err != nil {
			return nil, err
		}
		return inner.{{.Name}}(ctx, msg)
	}
	return msg, fn, nil
{{end -}}
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
{{end}}
`
