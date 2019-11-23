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
	"getRole":                  getRole,
	"getRoleLevel":             getRoleLevel,
	"requireSecret":            requireSecret,
	"getClaim":                 getClaim,
	"requireCaptcha":           requireCaptcha,
	"authPartner":              authPartner,
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
	return p.Type == permission.CurUsr
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

func getRole(m *Method) string {
	p := getPermission(m)
	return string(p.Role)
}

func getRoleLevel(m *Method) int {
	p := getPermission(m)
	lvl := permission.RoleLevel(p.Role)
	if p.Role != "" && lvl <= 0 {
		gen.Fatalf("Method %v: invalid role (%v)", m.FullPath(), p.Role)
	}
	return lvl
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
	case permission.Public, permission.Protected, permission.Custom, permission.Secret, permission.SuperAdmin:
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
type {{$m|methodName}}Endpoint struct {
	{{$m.Req|type}}
	Result {{$m.Resp|type}}
	Context {{$m|getClaim}}
	{{if $m|authPartner}}CtxPartner *model.Partner{{end}}
}
{{end}}

func (s wrap{{$s.Name}}Service) {{$m.Name}}(ctx context.Context, req {{.Req|type}}) (resp {{.Resp|type}}, err error) {
	t0 := time.Now()
{{if requireAuth $m -}}
	var session *middleware.Session
{{end -}}
	var errs []*cm.Error
	const rpcName = "{{$.PackageName}}.{{$s.Name}}/{{$m|methodName}}"
	defer func() {
		recovered := recover()
		err = cmwrapper.RecoverAndLog(ctx, rpcName, {{if requireAuth $m}}session{{else}}nil{{end}}, req, resp, recovered, err, errs, t0)
		metrics.CountRequest(rpcName, err)
	}()
	defer cmwrapper.Censor(req)
{{- if requireAuth $m}}
	sessionQuery := &middleware.StartSessionQuery{
		Context: ctx,
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
	if err := bus.Dispatch(ctx, sessionQuery); err != nil {
		return nil, err
	}
	session = sessionQuery.Result
{{- end}}
	query := &{{$m|methodName}}Endpoint{ {{$m.Req|baseName}}: req }
	{{if requireLogin $m -}}
	query.Context.Claim = session.Claim
	{{end -}}
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
	{{if authPartner $m | eq 0 | and (requireUser $m) -}}
	// Verify that the user has correct service type
	if session.Claim.AuthPartnerID != 0 {
		return nil, common.ErrPermissionDenied
	}
	{{end -}}
	{{if requireRole $m -}}
	query.Context.IsOwner = session.IsOwner
	query.Context.Roles = session.Roles
	query.Context.Permissions = session.Permissions
		{{if $m|getRole|ne "" -}}
	// Verify that the user has role "{{$m|getRole}}"
	if !session.IsOwner && permission.MaxRoleLevel(session.Roles) < {{$m|getRoleLevel}} {
		return nil, common.ErrPermissionDenied
	}
		{{end -}}
	{{end -}}
	{{if requireSecret $m -}}
	// Verify secret token
	token := middleware.GetBearerTokenFromCtx(ctx)
	if token != s.secret {
		return nil, common.ErrUnauthenticated
	}
	{{end -}}
	{{if requireCaptcha $m -}}
	// Verify captcha token
	if err := middleware.VerifyCaptcha(ctx, req.RecaptchaToken); err != nil {
		return nil, err
	}
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
	printer.Import("cm", "etop.vn/api/pb/common")
	printer.Import("cmwrapper", "etop.vn/backend/pkg/common/wrapper")
	printer.Import("common", "etop.vn/backend/pkg/common")
	printer.Import("metrics", "etop.vn/backend/pkg/common/metrics")
	printer.Import("middleware", "etop.vn/backend/pkg/etop/authorize/middleware")
	printer.Import("model", "etop.vn/backend/pkg/etop/model")
	printer.Import("permission", "etop.vn/backend/pkg/etop/authorize/permission")

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
