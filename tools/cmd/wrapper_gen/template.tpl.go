package main

const tplStr = `// package {{.PackageName}} generated by wrapper_gen. DO NOT EDIT.
package {{.PackageName}}W

import (
	"context"

	twirp "github.com/twitchtv/twirp"

	bus "etop.vn/backend/pkg/common/bus"
	common "etop.vn/backend/pkg/common"
	cmgrpc "etop.vn/backend/pkg/common/grpc"
	metrics "etop.vn/backend/pkg/common/metrics"
	telebot "etop.vn/backend/pkg/common/telebot"
	claims "etop.vn/backend/pkg/etop/authorize/claims"
	middleware "etop.vn/backend/pkg/etop/authorize/middleware"
	permission "etop.vn/backend/pkg/etop/authorize/permission"
	model "etop.vn/backend/pkg/etop/model"
	cm "etop.vn/backend/pb/common"
	cmwrapper "etop.vn/backend/pkg/common/wrapper"
	l "etop.vn/common/l"
	{{if .PackageName |eq "affiliate" -}}
	{{if .PackagePath |eq "etop.vn/backend/pb/etop/affiliate" -}}
	{{.PackageName}}z "etop.vn/backend/zexp/api/root/int/etop_{{.PackageName}}"
	{{else -}}
	{{.PackageName}}z "etop.vn/backend/zexp/api/root/int/services_{{.PackageName}}"
	{{end -}}
	{{else -}}
	{{if .PackageName |eq "shop" -}}
	{{if .PackagePath |eq "etop.vn/backend/pb/etop/shop" -}}
		{{.PackageName}}z "etop.vn/backend/zexp/api/root/int/{{.PackageName}}"
	{{else -}}
		{{.PackageName}}z "etop.vn/backend/zexp/api/root/int/external_{{.PackageName}}"
	{{end -}}
	{{else -}}
	{{.PackageName}}z "etop.vn/backend/zexp/api/root/int/{{.PackageName}}"
	{{end -}}
	{{end -}}
	{{.PackageName}} "{{.PackagePath}}"
	{{range .Imports -}}
	{{.Name}} "{{.Path}}"
	{{end}}
)

var ll = l.New()

type (
	EmptyClaim    		= claims.EmptyClaim
	UserClaim  	  		= claims.UserClaim
	AdminClaim    		= claims.AdminClaim
	PartnerClaim  		= claims.PartnerClaim
	ShopClaim     		= claims.ShopClaim
	AffiliateClaim    = claims.AffiliateClaim
)

type Muxer interface {
	Handle(string, http.Handler)
}

func New{{.ServiceName}}Server(mux Muxer, hooks *twirp.ServerHooks{{if $.HasSecret}}, secret string{{end}}) {
	{{- if $.HasSecret -}}
	if secret == "" {
		ll.Fatal("Secret is empty")
	}
	{{end -}}

{{range $s := .Services -}}
{{range $m := .Methods -}}
	bus.Expect(&{{$m|methodName}}Endpoint{})
{{end -}}
{{end -}}

{{range $s := .Services -}}
		mux.Handle({{$.PackageName}}z.{{$s.Name}}PathPrefix, {{$.PackageName}}z.New{{$s.Name}}Server({{$s.Name}}{ {{if $.HasSecret}}secret: secret{{end}} }))
{{end -}}
}

type {{.ServiceName}}Impl struct {
{{range $s := .Services -}}
	{{$s.Name}}
{{end -}}
}

{{range $s := .Services}}
type {{$s.Name}} struct { {{if $.HasSecret}}secret string{{end}} }

{{range $m := $s.Methods}}
type {{$m|methodName}}Endpoint struct {
	{{$m.InputType}}
	Result {{$m.OutputType}}
	Context {{claim $m}}
	{{if authPartner $m}}CtxPartner *model.Partner{{end}}
}

func (s {{$s.Name}}) {{$m.Name}}(ctx context.Context, req {{$m.InputType}}) (resp {{$m.OutputType}}, err error) {
	t0 := time.Now()
	{{if requireAuth $m -}}
	var session *middleware.Session
	{{end -}}
	var errs []*cm.Error
	const rpcName = "{{$.PackageName}}.{{trimName $s.Name}}/{{$m|methodName}}"
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
	query := &{{$m|methodName}}Endpoint{ {{baseName $m.InputType}}: req }
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
		{{if role $m | ne "" -}}
	// Verify that the user has role "{{role $m}}"
	if !session.IsOwner && permission.MaxRoleLevel(session.Roles) < {{roleLevel $m}} {
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
	err = bus.Dispatch(ctx, query)
	resp = query.Result
	if err == nil {
		if resp == nil {
			return nil, common.Error(common.Internal, "", nil).Log("nil response")
		}
		errs = cmwrapper.HasErrors(resp)
	}
	return resp, err
}
{{end}}
{{end}}
`
