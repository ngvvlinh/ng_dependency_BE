package auth

import (
	"fmt"
	"sort"
	"strings"

	"github.com/casbin/casbin"

	"o.o/backend/pkg/common/authorization/casbin/adapter"
)

type Policy string

const model = `
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

const Authorization = "e-authorization"

type Roles []string

type Authorizer struct {
	*casbin.Enforcer

	mapRoleAndActions map[string][]string
}

func New(policy Policy) *Authorizer {
	sa := adapter.NewAdapter(string(policy))
	enforcer := casbin.NewEnforcer(casbin.NewModel(model), sa)
	authorization := &Authorizer{
		Enforcer:          enforcer,
		mapRoleAndActions: buildMapRoleActions(string(policy)),
	}
	return authorization
}

type AuthorizationService interface {
	Check(roles []string, action string) bool
}

func (a *Authorizer) Check(roles []string, actionsArgs string, isTest int) bool {
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

func (a *Authorizer) CheckSingle(roles []string, action string) bool {
	for _, role := range roles {
		for _, act := range a.mapRoleAndActions[role] {
			if act == action {
				return true
			}
		}
	}
	return false
}

func (a *Authorizer) ListActionsByRolesMap(roles []string) map[string]struct{} {
	m := make(map[string]struct{})
	for _, role := range roles {
		for _, action := range a.mapRoleAndActions[role] {
			m[action] = struct{}{}
		}
	}
	return m
}

func (a *Authorizer) ListActionsByRoles(roles []string) (actions []string) {
	m := a.ListActionsByRolesMap(roles)
	for action := range m {
		actions = append(actions, action)
	}
	sort.Strings(actions)
	return
}

func buildMapRoleActions(policy string) map[string][]string {
	lines := strings.Split(policy, "\n")
	checkedActions := map[string]struct{}{}
	m := make(map[string][]string) // map role and actions
	for _, line := range lines {
		// prefix '#' for comment
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// p, readAndWriteBook, me, you => ["p",    "readAndWriteBook", "me", "you"]
		//                                   ^              ^                ^
		//                                 prefix         action          roles......
		elements := strings.Split(line, ",")
		if len(elements) < 3 {
			panic(fmt.Sprintf("Invalid policy setup, error line content: %v", line))
		}

		action := strings.TrimSpace(elements[1])
		if _, ok := checkedActions[action]; ok {
			panic(fmt.Sprintf("Duplicate action, error line content: %v", line))
		}
		checkedActions[action] = struct{}{}

		roles := elements[2:]
		for _, role := range roles {
			role = strings.TrimSpace(role)
			if role == "" {
				panic(fmt.Sprintf("Invalid policy setup, error line content: %v", line))
			}
			m[role] = append(m[role], action)
		}
	}
	return m
}
