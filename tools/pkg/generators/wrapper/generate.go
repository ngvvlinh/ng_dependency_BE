package wrapper

import (
	"go/types"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/pkg/etop/authorize/permission"
	"etop.vn/backend/tools/pkg/acl"
	"etop.vn/backend/tools/pkg/gen"
	"etop.vn/backend/tools/pkg/generator"
)

var currentPrinter generator.Printer
var funcs = map[string]interface{}{
	"type":                     renderType,
	"baseName":                 baseName,
	"hasSecret":                hasSecret,
	"isPublic":                 isPublic,
	"methodName":               methodName,
	"requireAuth":              requireAuth,
	"requireLogin":             requireLogin,
	"requireUser":              requireUser,
	"requireAPIKey":            requireAPIKey,
	"requireAPIPartnerShopKey": requireAPIPartnerShopKey,
	"requirePartner":           requirePartner,
	"requireShop":              requireShop,
	"requireAffiliate":         requireAffiliate,
	"requireEtopAdmin":         requireEtopAdmin,
	"requireSuperAdmin":        requireSuperAdmin,
	"requireRole":              requireRole,
	"requireSecret":            requireSecret,
	"getClaim":                 getClaim,
	"requireCaptcha":           requireCaptcha,
	"authPartner":              authPartner,
	"requireActions":           requireActions,
}

func getPermission(m *Method) *permission.PermissionDecl {
	p := acl.ACL[acl.ConvertKey(m.FullPath())]
	if p == nil {
		gen.Fatalf("Method %v does not exist in ACL map", m.FullPath())
	}
	return p
}

func renderType(t types.Type) string {
	return currentPrinter.TypeString(t)
}

func baseName(t types.Type) string {
	return t.(*types.Pointer).Elem().(*types.Named).Obj().Name()
}

var cacheSecret = map[string]bool{}

func hasSecret(s *Service) (result bool) {
	prefix := s.PkgName + "." + s.Name
	if s.PkgPrefix != "" {
		prefix = s.PkgPrefix + "/" + prefix
	}
	if has, ok := cacheSecret[prefix]; ok {
		return has
	}

	defer func() { cacheSecret[prefix] = result }()
	for k, perm := range acl.ACL {
		if strings.HasPrefix(k, prefix) && perm.Type == permission.Secret {
			return true
		}
	}
	return false
}

func methodName(m *Method) string {
	p := getPermission(m)
	if p.Rename != "" {
		return p.Rename
	}
	return m.Name
}

func isPublic(m *Method) bool {
	p := getPermission(m)
	return p.Type == permission.Public
}

func requireAuth(m *Method) bool {
	p := getPermission(m)
	return p.Type >= permission.Protected
}

func requireLogin(m *Method) bool {
	p := getPermission(m)
	return p.Type >= permission.Protected && p.Type <= permission.EtopAdmin
}

func requireUser(m *Method) bool {
	p := getPermission(m)
	return p.Type == permission.CurUsr || p.Auth == permission.User
}

func requireAPIKey(m *Method) bool {
	p := getPermission(m)
	return p.Auth == permission.APIKey
}

func requireAPIPartnerShopKey(m *Method) bool {
	p := getPermission(m)
	return p.Auth == permission.APIPartnerShopKey
}

func requirePartner(m *Method) bool {
	p := getPermission(m)
	return p.Type == permission.Partner
}

func requireShop(m *Method) bool {
	p := getPermission(m)
	return p.Type == permission.Shop
}

func requireAffiliate(m *Method) bool {
	p := getPermission(m)
	return p.Type == permission.Affiliate
}

func requireEtopAdmin(m *Method) bool {
	p := getPermission(m)
	return p.Type == permission.EtopAdmin
}

func requireSuperAdmin(m *Method) bool {
	p := getPermission(m)
	return p.Type == permission.SuperAdmin
}

func requireRole(m *Method) bool {
	p := getPermission(m)
	return p.Type >= permission.Partner && p.Type <= permission.EtopAdmin
}

func requireSecret(m *Method) bool {
	p := getPermission(m)
	return p.Type == permission.Secret
}

func getClaim(m *Method) string {
	p := getPermission(m)
	switch p.Type {
	case permission.CurUsr:
		return "claims.UserClaim"
	case permission.EtopAdmin:
		return "claims.AdminClaim"
	case permission.Partner:
		return "claims.PartnerClaim"
	case permission.Shop:
		return "claims.ShopClaim"
	case permission.Affiliate:
		return "claims.AffiliateClaim"
	case permission.Public, permission.Protected, permission.Secret, permission.SuperAdmin:
		return "claims.EmptyClaim"
	}

	gen.Fatalf("Method %v: invalid permission (%v)", m.FullPath(), p.Type)
	return ""
}

func requireCaptcha(m *Method) string {
	p := getPermission(m)
	return p.Captcha
}

func authPartner(m *Method) int {
	p := getPermission(m)
	return int(p.AuthPartner)
}

func requireActions(m *Method) string {
	p := getPermission(m)
	var actionsTemp []string
	for _, action := range p.Actions {
		actionsTemp = append(actionsTemp, string(action))
	}
	return strings.Join(actionsTemp, "|")
}

var tpl = template.Must(template.New("tpl").Funcs(funcs).Parse(tplText))

const tplText = `
{{range $s := .Services}}
func Wrap{{.Name}}Service(s *{{.Name}}Service{{if $s|hasSecret}}, secret string{{end}}) api.{{.Name}}Service {
	return wrap{{.Name}}Service{s: s{{if $s|hasSecret}}, secret: secret{{end}}}
}

type wrap{{.Name}}Service struct {
	s *{{.Name}}Service
{{if $s|hasSecret}}secret string{{end}}
}

{{range $m := .Methods}}
{{if $m.Kind|eq 1}}
type {{$s.EndpointPrefix}}{{$m|methodName}}Endpoint struct {
	{{$m.Req|type}}
	Result {{$m.Resp|type}}
	Context {{$m|getClaim}}
	{{if $m|authPartner}}CtxPartner *identitymodel.Partner{{end}}
}
{{end}}

func (s wrap{{$s.Name}}Service) {{$m.Name}}(ctx context.Context, req {{.Req|type}}) (resp {{.Resp|type}}, err error) {
	t0 := time.Now()
	var session *middleware.Session
	var errs []*cm.Error
	const rpcName = "{{$.PackageName}}.{{$s.Name}}/{{$m|methodName}}"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, {{if requireAuth $m}}session{{else}}nil{{end}}, req, resp, recovered, err, errs, t0)
	}()
	defer cmwrapper.Censor(req)
	sessionQuery := &middleware.StartSessionQuery{
		{{if requireAuth       $m}}RequireAuth: 	  true,
{{end}}	{{if requireUser       $m}}RequireUser: 	  true,
{{end}} {{if requireAPIKey     $m}}RequireAPIKey:     true,
{{end}} {{if requireAPIPartnerShopKey     $m}}RequireAPIPartnerShopKey:     true,
{{end}} {{if requirePartner    $m}}RequirePartner:    true,
{{end}}	{{if requireShop       $m}}RequireShop:       true,
{{end}} {{if requireAffiliate  $m}}RequireAffiliate: 			true,
{{end}}	{{if requireEtopAdmin  $m}}RequireEtopAdmin:  true,
{{end}}	{{if requireSuperAdmin $m}}RequireSuperAdmin: true,
{{end}} {{if authPartner       $m}}AuthPartner: {{authPartner $m}},
{{end}}
	}
	ctx, err = middleware.StartSession(ctx, sessionQuery)
	if err != nil {
		{{if isPublic $m -}}
		// ignore invalid authentication token
		if common.ErrorCode(err) != common.Unauthenticated {
			return nil, err
		}
		{{else -}}
		return nil, err
		{{end -}}
	}
	session = sessionQuery.Result
	query := &{{$s.EndpointPrefix}}{{$m|methodName}}Endpoint{ {{$m.Req|baseName}}: req }
	if session != nil {
		query.Context.Claim = session.Claim
	}
	{{if requireUser $m -}}
	query.Context.User = session.User
	query.Context.Admin = session.Admin
	{{end -}}
	{{if requirePartner $m -}}
	query.Context.Partner = session.Partner
	{{end -}}
	{{if requireShop $m -}}
	query.Context.Shop = session.Shop
	{{end -}}
	{{if requireAffiliate $m -}}
	query.Context.Affiliate = session.Affiliate
	{{end -}}
	{{if requireEtopAdmin $m -}}
	query.Context.IsEtopAdmin = session.IsEtopAdmin
	{{end -}}
	{{if requireSuperAdmin $m -}}
	query.Context.IsSuperAdmin = session.IsSuperAdmin
	{{end -}}
	{{if authPartner $m -}}
	query.CtxPartner = session.CtxPartner
	{{end -}}
	{{if requireRole $m -}}
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
	{{end -}}

	{{if requireActions $m -}}
	isTest := 0
	if query.Context.Shop != nil {
		isTest = query.Context.Shop.IsTest
	}
	authorization := auth.New()
	// Do not check permission for 3rd party requests
	if session.Claim.AuthPartnerID == 0 && !authorization.Check(query.Context.Roles, "{{$m|requireActions}}", isTest) {
		return nil, common.Error(common.PermissionDenied, "", nil)
	}
	query.Context.Actions = strings.Split("{{$m|requireActions}}", "|")
	{{end -}}

	{{if requireSecret $m -}}
	// Verify secret token
	token := headers.GetBearerTokenFromCtx(ctx)
	if token != s.secret {
		return nil, common.ErrUnauthenticated
	}
	{{end -}}
	{{if requireCaptcha $m -}}
	// Verify captcha token
	{{if requireCaptcha $m | eq "custom" -}}
	if req.RecaptchaToken != "" || req.RequireCaptcha() {
	{{end -}}
	if err := middleware.VerifyCaptcha(ctx, req.RecaptchaToken); err != nil {
		return nil, err
	}
	{{if requireCaptcha $m | eq "custom" -}}
	}
	{{end -}}
	{{end -}}
	ctx = bus.NewRootContext(ctx)
	{{if $m.Kind|eq 1 -}}
	err = s.s.{{$m.Name}}(ctx, query)
	resp = query.Result
	{{else -}}
	// TODO
	resp, err = s.s.{{$m.Name}}(ctx, claim, req)
	{{end -}}
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, common.Error(common.Internal, "", nil).Log("nil response")
	}
	errs = cmwrapper.HasErrors(resp)
	return resp, nil 
}
{{end}}
{{end}}
`

func generate(printer generator.Printer, pkg *packages.Package, services []*Service) error {
	if len(services) == 0 {
		return nil
	}
	currentPrinter = printer
	printer.Import("api", services[0].PkgPb)
	printer.Import("bus", "etop.vn/backend/pkg/common/bus")
	printer.Import("claims", "etop.vn/backend/pkg/etop/authorize/claims")
	printer.Import("cm", "etop.vn/api/top/types/common")
	printer.Import("cmwrapper", "etop.vn/backend/pkg/common/apifw/wrapper")
	printer.Import("common", "etop.vn/backend/pkg/common")
	printer.Import("metrics", "etop.vn/backend/pkg/common/metrics")
	printer.Import("middleware", "etop.vn/backend/pkg/etop/authorize/middleware")
	printer.Import("identitymodel", "etop.vn/backend/com/main/identity/model")
	printer.Import("permission", "etop.vn/backend/pkg/etop/authorize/permission")
	printer.Import("auth", "etop.vn/backend/pkg/etop/authorize/auth")
	printer.Import("headers", "etop.vn/backend/pkg/common/headers")

	vars := map[string]interface{}{
		"PackageName":   services[0].PkgName,
		"PackagePrefix": services[0].PkgPrefix,
		"Services":      services,
	}
	err := tpl.Execute(printer, vars)
	if err != nil {
		return generator.Errorf(err, "generate %v: %v", pkg.Name, err)
	}
	return nil
}
