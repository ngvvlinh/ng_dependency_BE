package genutil

import (
	"go/types"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize/english"

	"etop.vn/backend/tools/pkg/generator"
)

var _ generator.Qualifier = &Qualifier{}

type Qualifier struct{}

func (q Qualifier) Qualify(pkg *types.Package) string {
	alias := pkg.Name()
	if alias == "model" || alias == "types" || alias == "convert" {
		super := filepath.Base(filepath.Dir(pkg.Path()))
		alias = strings.ToLower(super) + alias
	}
	return alias
}

func Plural(s string) string {
	return english.PluralWord(2, s, "")
}
