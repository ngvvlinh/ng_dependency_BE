package genutil

import (
	"go/types"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize/english"
	"golang.org/x/tools/go/packages"

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

func HasPrefixCamel(s string, prefix string) bool {
	ln := len(prefix)
	return len(s) > ln &&
		s[:ln] == prefix &&
		!(s[ln] >= 'a' && s[ln] <= 'z')
}

func NoUnboundDirectives(ng generator.Engine, pkg *packages.Package, prefix string) error {
	ds := ng.GetDirectivesByPackage(pkg).FilterBy("sqlgen")
	if len(ds) == 0 {
		return nil
	}
	var s strings.Builder
	for i := 0; i < len(ds) && i < 10; i++ {
		s.WriteString("\n\t")
		s.WriteString(ds[i].Raw)
	}
	return generator.Errorf(nil, "package %v: found unbound directives%v", pkg.PkgPath, s.String())
}
