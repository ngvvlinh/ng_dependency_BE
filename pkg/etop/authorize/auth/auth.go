package auth

import (
	"strings"

	"github.com/casbin/casbin"

	"o.o/backend/pkg/common/authorization/casbin/adapter"
)

const (
	model = `
	[request_definition]
	r = action, role
	
	[policy_definition]
	p = action, role
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = g(r.role, p.role) && r.action == p.action`
)

const (
	EcomAuthorization = "e-authorization"
)

var authoration *Authorization

func Init(policy string) {
	sa := adapter.NewAdapter(policy)
	enforcer := casbin.NewEnforcer(casbin.NewModel(model), sa)
	authoration = &Authorization{Enforcer: enforcer, Policy: policy}
}

type Authorization struct {
	*casbin.Enforcer
	Policy string
}

type AuthorizationService interface {
	Check(roles []string, action string) bool
}

func New() *Authorization {
	return authoration
}

func (a *Authorization) Check(roles []string, actionsArgs string, isTest int) bool {
	actions := strings.Split(actionsArgs, "|")
	for _, role := range roles {
		for _, action := range actions {
			if a.Enforcer.Enforce(action, role) {
				return true
			}
		}
	}
	return false
}

type Roles []string

func (r Roles) Check(actionArg string) bool {
	actions := ListActionsByRoles(r)

	for _, action := range actions {
		if actionArg == action {
			return true
		}
	}
	return false
}

func ListActionsByRoles(roles []string) (actions []string) {
	strs := strings.Split(authoration.Policy, "\n")

	mapRoleAndActions := make(map[string][]string)
	for _, str := range strs {
		// prefix '#' for comment
		str = strings.TrimSpace(str)
		if str == "" || strings.HasPrefix(str, "#") {
			continue
		}

		elements := strings.Split(str, ",")
		_ = elements[0] // prefix 'p' | 'g'
		action := strings.TrimSpace(elements[1])
		roles := elements[2:]
		for _, role := range roles {
			role = strings.TrimSpace(role)
			mapRoleAndActions[role] = append(mapRoleAndActions[role], action)
		}
	}

	mapAction := make(map[string]bool)
	for _, role := range roles {
		for _, action := range mapRoleAndActions[role] {
			mapAction[action] = true
		}
	}

	for action := range mapAction {
		actions = append(actions, action)
	}
	return
}
