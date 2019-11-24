package apix

const tplText = `
type Server interface {
	http.Handler
	PathPrefix() string
}

{{range $s := .Services}}

type {{.Name}}ServiceServer struct {
	inner {{.Name}}Service
}

func New{{.Name}}ServiceServer(svc {{.Name}}Service) Server{
	return &{{.Name}}ServiceServer {
		inner: svc,
	}
}

const {{.Name}}ServicePathPrefix = "{{.APIPath}}/"
func (s *{{$s.Name}}ServiceServer)PathPrefix() string {
	return {{.Name}}ServicePathPrefix
}

func (s *{{$s.Name}}ServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path)
	if err != nil {
		httprpc.WriteError(ctx, resp, err)
		return
	}
	serve(ctx, resp, req, reqMsg, exec)
}

func (s *{{$s.Name}}ServiceServer) parseRoute(path string) (reqMsg capi.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
{{range $m := .Methods -}}
	case "{{$s.APIPath}}/{{.Name}}":
	msg := {{(index .Request.Items 0).Type|new}}
	fn := func(ctx context.Context) (capi.Message, error) {
		return s.inner.{{.Name}}(ctx, msg)
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
