package sqlgen

import (
	"go/types"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"

	"o.o/backend/tools/pkg/generator"
	"o.o/backend/tools/pkg/generators/api/parse"
	"o.o/backend/tools/pkg/genutil"
	"o.o/common/l"
)

const CmdPrefix = "sqlgen"

var ll = l.New()
var CurrentInfo *parse.Info

type plugin struct {
	generator.Qualifier
}

func New() generator.Plugin {
	return &plugin{
		Qualifier: genutil.Qualifier{},
	}
}

func (p *plugin) Name() string { return CmdPrefix }

func (p *plugin) Filter(ng generator.FilterEngine) error {
	CurrentInfo = parse.NewInfo(ng)
	return generator.FilterByCommand(CmdPrefix).FilterAll(ng)
}

func (p *plugin) Generate(ng generator.Engine) error {
	CurrentInfo.Init(ng)
	return ng.GenerateEachPackage(p.generateEachPackage)
}

func (p *plugin) generateEachPackage(ng generator.Engine, pkg *packages.Package, printer generator.Printer) error {
	if err := genutil.NoUnboundDirectives(ng, pkg, CmdPrefix); err != nil {
		return err
	}

	pr = printer
	g := &genImpl{
		ng:      ng,
		Printer: printer,
		mapBase: make(map[string]bool),
		mapType: make(map[string]*typeDef),
	}
	var typs []*types.Named
	for _, obj := range ng.GetObjectsByPackage(pkg) {
		typ, err := parseDirectives(ng, pkg, obj, g)
		if err != nil {
			return err
		}
		if typ != nil {
			typs = append(typs, typ)
		}
	}

	printer.Import("cmsql", "o.o/backend/pkg/common/sql/cmsql")
	printer.Import("core", "o.o/backend/pkg/common/sql/sq/core")
	printer.Import("migration", "o.o/backend/pkg/common/sql/migration")

	sort.Slice(typs, func(i, j int) bool {
		return typs[i].Obj().Name() < typs[j].Obj().Name()
	})
	for _, typ := range typs {
		if err := g.Generate(typ); err != nil {
			return err
		}
	}
	return nil
}

//                                       1..........12........3..........32
var reFrom = regexp.MustCompile(`^\s*([A-z0-9_]+)(\s+as\s+([A-z0-9_]+))?$`)

//                                       1..........12........3..........32.........4..4
var reJoin = regexp.MustCompile(`^\s*([A-z0-9_]+)(\s+as\s+([A-z0-9_]+))?\s+on\s+(.+)$`)

func parseDirectives(ng generator.Engine, pkg *packages.Package, obj types.Object, g *genImpl) (*types.Named, error) {
	ds := ng.GetDirectives(obj).FilterBy(CmdPrefix)
	if len(ds) == 0 {
		return nil, nil
	}

	// must be struct type
	typeName, ok := obj.(*types.TypeName)
	if !ok {
		return nil, generator.Errorf(nil, "sqlgen must be used on struct declaration (%v is not)", obj)
	}
	named := typeName.Type().(*types.Named)
	if st := unwrapNamedStruct(named); st == nil {
		return nil, generator.Errorf(nil, "sqlgen must be used on struct declaration (%v is not)", obj)
	}

	var baseType *types.Named
	opts := []Option{OptionSimple(named)}
	for _, d := range ds {
		switch {
		case d.Cmd == CmdPrefix && d.Arg == "":
			opts = append(opts, OptionSimple(named))
			baseType = named

		case d.Cmd == CmdPrefix:
			baseName, alias, err := parseDirectiveFrom(d)
			if err != nil {
				return nil, err
			}

			// parse option: derived
			// heuristic:
			// - when join, use field name
			// - when not join, use scope name
			var baseNamedStruct *types.Named
			if !directivesContain(ds, "join") {
				baseNamedStruct = getNamedStruct(pkg, baseName)
			} else {
				baseNamedStruct, err = getNamedStructFromField(named, baseName)
			}
			if err != nil {
				return nil, generator.Errorf(err, "%v: %v (%v)", obj.Name(), err, d.Raw)
			}
			if baseNamedStruct == nil {
				return nil, generator.Errorf(nil, "%v: %v not found (%v)", obj.Name(), baseName, d.Raw)
			}

			opts = append(opts, OptionDerived(named, baseNamedStruct))
			baseType = baseNamedStruct

			// parse option: alias
			if alias != "" {
				opts = append(opts, OptionAs(alias))
			}

		case strings.HasSuffix(d.Cmd, "join"):
			joinKeyword := strings.TrimPrefix(d.Cmd, CmdPrefix+":")
			switch joinKeyword {
			case "join", "left-join", "right-join", "full-join":
				// continue
			default:
				return nil, generator.Errorf(nil, "invalid join directive (%v)", d.Raw)
			}
			joinKeyword = strings.ReplaceAll(joinKeyword, "-", " ")
			joinKeyword = strings.ToUpper(joinKeyword)

			// parse option: join
			joinName, joinAlias, joinCond, err := parseDirectiveJoin(d)
			if err != nil {
				return nil, err
			}
			joinNamedStruct, err := getNamedStructFromField(named, joinName)
			if err != nil {
				return nil, generator.Errorf(err, "%v: %v (%v)", obj.Name(), err, d.Raw)
			}
			if joinNamedStruct == nil {
				return nil, generator.Errorf(nil, "%v: %v not found (%v)", obj.Name(), joinName, d.Raw)
			}
			opts = append(opts, OptionJoin(baseType, joinNamedStruct, joinKeyword, joinAlias, joinCond))

		default:
			return nil, generator.Errorf(nil, "invalid directive (%v)", d.Raw)
		}
	}
	return named, g.AddStruct(named, opts...)
}

func directivesContain(ds generator.Directives, substr string) bool {
	for _, d := range ds {
		if strings.Contains(d.Cmd, substr) {
			return true
		}
	}
	return false
}

func parseDirectiveFrom(d generator.Directive) (baseName string, alias string, _ error) {
	parts := reFrom.FindStringSubmatch(d.Arg)
	if len(parts) == 0 {
		return "", "", generator.Errorf(nil, "sqlgen directive is invalid (%v)", d.Raw)
	}
	baseName = parts[1]
	if len(parts) >= 4 {
		alias = parts[3]
	}
	return baseName, alias, nil
}

func parseDirectiveJoin(d generator.Directive) (typ, alias, cond string, _ error) {
	parts := reJoin.FindStringSubmatch(d.Arg)
	if len(parts) == 0 {
		return "", "", "", generator.Errorf(nil, "invalid join directive description (%v)", d.Raw)
	}
	return parts[1], parts[3], parts[4], nil
}

func getNamedStruct(pkg *packages.Package, name string) *types.Named {
	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		return nil
	}
	typ, ok := obj.(*types.TypeName)
	if !ok {
		return nil
	}
	namedStruct := typ.Type().(*types.Named)
	if st := unwrapNamedStruct(namedStruct); st == nil {
		return nil
	}
	return namedStruct
}

func getNamedStructFromField(named *types.Named, fieldName string) (*types.Named, error) {
	st := unwrapNamedStruct(named)
	for i, n := 0, st.NumFields(); i < n; i++ {
		field := st.Field(i)
		if field.Name() == fieldName {
			fieldType, ok := pointerToNamedStruct(field.Type())
			if !ok {
				return nil, generator.Errorf(nil, "field %v.%v must be a pointer to named struct", named.Obj().Name(), fieldName)
			}
			return fieldType, nil
		}
	}
	return nil, generator.Errorf(nil, "field %v.%v not found", named.Obj().Name(), fieldName)
}

func unwrapNamedStruct(named *types.Named) *types.Struct {
	underlying := genutil.UnwrapNamed(named)
	st, _ := underlying.(*types.Struct)
	return st
}
