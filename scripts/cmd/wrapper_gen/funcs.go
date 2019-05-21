package main

import (
	"regexp"
	"text/template"

	"etop.vn/backend/pkg/etop/authorize/permission"
	"etop.vn/backend/up/gogen/pkg/gen"
	g "etop.vn/backend/up/gogen/pkg/grpcgen"
)

var funcs = template.FuncMap{
	"baseName":                 baseName,
	"trimName":                 trimName,
	"requireAuth":              requireAuth,
	"requireLogin":             requireLogin,
	"requireUser":              requireUser,
	"requirePartner":           requirePartner,
	"requireAPIKey":            requireAPIKey,
	"requireAPIPartnerShopKey": requireAPIPartnerShopKey,
	"requireShop":              requireShop,
	"requireEtopAdmin":         requireEtopAdmin,
	"requireSuperAdmin":        requireSuperAdmin,
	"requireRole":              requireRole,
	"requireSecret":            requireSecret,
	"role":                     getRole,
	"roleLevel":                getRoleLevel,
	"claim":                    getClaim,
	"requireCaptcha":           requireCaptcha,
	"authPartner":              authPartner,
}

func getPermission(m *g.Method) *permission.PermissionDecl {
	p := ACL[m.FullPath()]
	if p == nil {
		gen.Fatalf("Method %v does not exist in ACL map", m.FullPath())
	}
	return p
}

var baseNameRegex = regexp.MustCompile(`[0-9A-z]+$`)

func baseName(name string) string {
	return baseNameRegex.FindString(name)
}

func trimName(name string) string {
	return name[:len(name)-len("Service")]
}

func requireAuth(m *g.Method) bool {
	p := getPermission(m)
	return p.Type >= Protected
}

func requireLogin(m *g.Method) bool {
	p := getPermission(m)
	return p.Type >= Protected && p.Type <= EtopAdmin
}

func requireUser(m *g.Method) bool {
	p := getPermission(m)
	return p.Type == CurUsr
}

func requireAPIKey(m *g.Method) bool {
	p := getPermission(m)
	return p.Auth == permission.APIKey
}

func requireAPIPartnerShopKey(m *g.Method) bool {
	p := getPermission(m)
	return p.Auth == permission.APIPartnerShopKey
}

func requirePartner(m *g.Method) bool {
	p := getPermission(m)
	return p.Type == Partner
}

func requireShop(m *g.Method) bool {
	p := getPermission(m)
	return p.Type == Shop
}

func requireEtopAdmin(m *g.Method) bool {
	p := getPermission(m)
	return p.Type == EtopAdmin
}

func requireSuperAdmin(m *g.Method) bool {
	p := getPermission(m)
	return p.Type == SuperAdmin
}

func requireRole(m *g.Method) bool {
	p := getPermission(m)
	return p.Type >= Partner && p.Type <= EtopAdmin
}

func getRole(m *g.Method) string {
	p := getPermission(m)
	return string(p.Role)
}

func getRoleLevel(m *g.Method) int {
	p := getPermission(m)
	lvl := permission.RoleLevel(p.Role)
	if p.Role != "" && lvl <= 0 {
		gen.Fatalf("Method %v: invalid role (%v)", m.FullPath(), p.Role)
	}
	return lvl
}

func requireSecret(m *g.Method) bool {
	p := getPermission(m)
	return p.Type == Secret
}

func getClaim(m *g.Method) string {
	p := getPermission(m)
	switch p.Type {
	case CurUsr:
		return "UserClaim"
	case EtopAdmin:
		return "AdminClaim"
	case Partner:
		return "PartnerClaim"
	case Shop:
		return "ShopClaim"
	case Supplier:
		return "SupplierClaim"
	case Public, Protected, Custom, Secret, SuperAdmin:
		return "EmptyClaim"
	}

	gen.Fatalf("Method %v: invalid permission (%v)", m.FullPath(), p.Type)
	return ""
}

func requireCaptcha(m *g.Method) string {
	p := getPermission(m)
	return p.Captcha
}

func authPartner(m *g.Method) int {
	p := getPermission(m)
	return int(p.AuthPartner)
}
