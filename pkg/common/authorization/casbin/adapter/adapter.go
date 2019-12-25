package adapter

import (
	"errors"
	"strings"

	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
)

type Adapter struct {
	Line string
}

func NewAdapter(line string) *Adapter {
	return &Adapter{
		Line: line,
	}
}

func (sa *Adapter) LoadPolicy(model model.Model) error {
	if sa.Line == "" {
		return errors.New("invalid line, line cannot be empty")
	}
	strs := strings.Split(sa.Line, "\n")

	for _, str := range strs {
		// prefix '#' for comment
		str = strings.TrimSpace(str)
		if str == "" || strings.HasPrefix(str, "#") {
			continue
		}

		elements := strings.Split(str, ",")
		prefix := strings.TrimSpace(elements[0])
		action := strings.TrimSpace(elements[1])
		roles := elements[2:]

		for _, role := range roles {
			role = strings.TrimSpace(role)
			persist.LoadPolicyLine(strings.Join([]string{prefix, action, role}, ", "), model)
		}
	}

	return nil
}

func (sa *Adapter) SavePolicy(model model.Model) error {
	panic("implement me")
}

func (sa *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	panic("implement me")
}

func (sa *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	panic("implement me")
}

func (sa *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	panic("implement me")
}
