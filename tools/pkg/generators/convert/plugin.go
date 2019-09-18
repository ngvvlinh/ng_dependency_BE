package convert

import (
	"fmt"
	"go/types"
	"io"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/genutil"
	"etop.vn/common/l"
)

var ll = l.New()

const Command = "gen:convert"
const ModeType = "convert:type"
const ModeCreate = "convert:create"
const ModeUpdate = "convert:update"

func New() generator.Plugin {
	return &plugin{
		Filter:    generator.FilterByCommand(Command),
		Qualifier: genutil.Qualifier{},
	}
}

type plugin struct {
	generator.Filter
	generator.Qualifier
}

func (p *plugin) Name() string { return "convert" }

func (p *plugin) Generate(ng generator.Engine) error {
	pkgs := ng.GeneratingPackages()
	for _, pkg := range pkgs {
		err := generatePackage(ng, pkg)
		if err != nil {
			return err
		}
	}
	return nil
}

type objMap struct {
	src  generator.Object
	gens []objGen
}

type objGen struct {
	mode string
	obj  generator.Object
	opts options
}

type objName struct {
	pkg  string
	name string
}

type options struct {
	identifiers []string
}

type fieldConvert struct {
	Out          *types.Var
	Arg          *types.Var
	IsIdentifier bool // for updating only
}

func (o objName) String() string {
	if o.pkg == "" {
		return o.name
	}
	return o.pkg + "." + o.name
}

func generatePackage(ng generator.Engine, gpkg *generator.GeneratingPackage) error {
	count := 0
	for _, d := range gpkg.Directives {
		if d.Cmd != Command {
			continue
		}
		count++
		apiPkgPaths, toPkgPaths, err := parseConvertDirective(d)
		if err != nil {
			return err
		}
		currentPrinter = gpkg.Generate()
		err = generatePackageStep(ng, currentPrinter, apiPkgPaths, toPkgPaths)
		if err != nil {
			return err
		}
	}
	if count == 0 {
		return generator.Errorf(nil, "invalid directive (must in format pkg1 -> pkg2)")
	}
	return nil
}

func generatePackageStep(ng generator.Engine, p generator.Printer, apiPkgPaths, toPkgPaths []string) error {
	ll.V(1).Debugf("convert from %v to %v", strings.Join(apiPkgPaths, ","), strings.Join(toPkgPaths, ","))

	flagSelf, err := validateEquality(apiPkgPaths, toPkgPaths)
	if err != nil {
		return err
	}
	flagAuto := !flagSelf && len(apiPkgPaths) == 1 && len(toPkgPaths) == 1
	flagInfo := false

	apiPkgs := make([]*packages.Package, len(apiPkgPaths))
	for i, pkgPath := range apiPkgPaths {
		apiPkgs[i] = ng.PackageByPath(pkgPath)
		if apiPkgs[i] == nil {
			return generator.Errorf(nil, "can not find package %v", pkgPath)
		}
	}
	toPkgs := make([]*packages.Package, len(toPkgPaths))
	for i, pkgPath := range toPkgPaths {
		toPkgs[i] = ng.PackageByPath(pkgPath)
		if toPkgs[i] == nil {
			return generator.Errorf(nil, "can not find package %v", pkgPath)
		}
	}

	apiObjMap := make(map[objName]*objMap)
	for _, pkg := range apiPkgs {
		apiObjs := ng.ObjectsByPackage(pkg)
		for _, obj := range apiObjs {
			if !obj.Object.Exported() {
				continue
			}
			objName := objName{pkg: pkg.PkgPath, name: obj.Object.Name()}
			apiObjMap[objName] = &objMap{src: obj}
		}
	}

	if ll.Verbosed(3) {
		for objName, objMap := range apiObjMap {
			ll.V(3).Debugf("object %v: %v", objName, objMap.src.Object.Type())
		}
	}

	for _, toPkg := range toPkgs {
		toObjs := ng.ObjectsByPackage(toPkg)
		for _, obj := range toObjs {
			ll.V(2).Debugf("convert to object %v with directives %#v", obj.Object.Name(), obj.Directives)
			if !obj.Object.Exported() {
				continue
			}
			if s := validateStruct(obj.Object); s == nil {
				continue
			}
			raw, mode, name, opts, err := parseWithMode(apiPkgs, obj.Directives)
			if err != nil {
				return err
			}
			ll.V(3).Debugf("parsed type %v with mode %v", name, mode)
			if mode == "" {
				// automatically convert type with the same name
				if flagAuto {
					name = objName{apiPkgs[0].PkgPath, obj.Object.Name()}
					if apiObjMap[name] == nil {
						continue
					}
					mode = ModeType

				} else {
					// notify user that there is no type generated for the conversion
					flagInfo = true
				}

			} else if apiObjMap[name] == nil {
				return generator.Errorf(nil, "type %v not found (directive %v)", name, raw)
			}
			if (name == objName{}) {
				continue
			}

			m := apiObjMap[name]
			if s := validateStruct(m.src.Object); s == nil {
				return generator.Errorf(nil, "%v is not a struct", m.src.Object.Name())
			}
			m.gens = append(m.gens, objGen{mode: mode, obj: obj, opts: opts})
		}
	}

	count, err := generateConverts(p, apiObjMap)
	if err != nil {
		return err
	}
	if flagInfo && count == 0 {
		ll.Warn("no types generated (for multiple package conversion, must use convert:type to define mapping)")
	}
	return nil
}

func validateEquality(pkgs1, pkgs2 []string) (bool, error) {
	if len(pkgs1) != len(pkgs2) {
		return false, nil
	}
	count := 0
	flags := make([]bool, len(pkgs1))
	for i, p1 := range pkgs1 {
		for _, p2 := range pkgs2 {
			if p1 == p2 {
				if flags[i] {
					return false, generator.Errorf(nil, "duplicated package (%v)", p1)
				}
				count++
				flags[i] = true
			}
		}
	}
	return count == len(pkgs1), nil
}

func validateStruct(obj types.Object) *types.Struct {
	s, ok := obj.(*types.TypeName)
	if !ok {
		return nil
	}
	st, _ := s.Type().Underlying().(*types.Struct)
	return st
}

func parseConvertDirective(directive generator.Directive) (apiPkgs, toPkgs []string, err error) {
	// parse "pkg" without "->"
	if strings.Index(directive.Arg, "->") < 0 {
		pkgs := strings.Split(directive.Arg, ",")
		for i := range pkgs {
			pkgs[i] = strings.TrimSpace(pkgs[i])
		}
		return pkgs, pkgs, nil
	}

	// parse "pkg1 -> pkg2"
	parts := strings.Split(directive.Arg, "->")
	if len(parts) != 2 {
		err = generator.Errorf(nil, "invalid directive (must in format pkg1 -> pkg2)")
		return
	}
	// parse "pkg1,pkg2"
	toPkgs = strings.Split(parts[0], ",")
	for i := range toPkgs {
		toPkgs[i] = strings.TrimSpace(toPkgs[i])
		if toPkgs[i] == "" {
			err = generator.Errorf(nil, "invalid directive (must in format pkg1 -> pkg2)")
			return
		}
	}
	// parse "pkg1,pkg2"
	apiPkgs = strings.Split(parts[1], ",")
	for i := range apiPkgs {
		apiPkgs[i] = strings.TrimSpace(apiPkgs[i])
		if apiPkgs[i] == "" {
			err = generator.Errorf(nil, "invalid directive (must in format pkg1 -> pkg2)")
			return
		}
	}
	return apiPkgs, toPkgs, nil
}

var reName = regexp.MustCompile(`[A-Z][A-z0-9_]*`)

func parseWithMode(apiPkgs []*packages.Package, directives []generator.Directive) (raw, mode string, _ objName, _ options, _ error) {
	for _, d := range directives {
		switch d.Cmd {
		case ModeType, ModeCreate, ModeUpdate:
			objName, opts, err := parseTypeName(apiPkgs, d.Arg)
			if err == nil && d.Cmd != ModeUpdate {
				if len(opts.identifiers) != 0 {
					err = generator.Errorf(nil, "invalid extra option (%v)", d.Arg)
				}
			}
			return d.Raw, d.Cmd, objName, opts, err
		}
	}
	return
}

var reTypeName = regexp.MustCompile(`^(.+\.)?([^.(]+)(\([^)]*\))?$`)

func parseTypeName(apiPkgs []*packages.Package, input string) (_ objName, opts options, err error) {
	parts := reTypeName.FindStringSubmatch(input)
	if len(parts) == 0 {
		err = generator.Errorf(nil, "invalid convert directive (%v)", input)
		return
	}
	pkgPath, name, extra := parts[1], parts[2], parts[3]
	if pkgPath != "" {
		pkgPath = pkgPath[:len(pkgPath)-1] // remove "."
	}
	if extra != "" {
		extra = extra[1 : len(extra)-1] // remove "(" ")"
		opts.identifiers = strings.Split(extra, ",")
		for _, ident := range opts.identifiers {
			if !reName.MatchString(ident) {
				err = generator.Errorf(nil, "invalid field name (%v)", input)
				return
			}
		}
	}

	if pkgPath == "" {
		if !reName.MatchString(name) {
			err = generator.Errorf(nil, "invalid type name (%v)", input)
			return
		}
		if len(apiPkgs) == 1 {
			return objName{pkg: apiPkgs[0].PkgPath, name: name}, opts, nil
		}
		err = generator.Errorf(nil, "must provide path for multiple input packages (%v)", input)
		return
	}

	pkgPaths := make([]string, len(apiPkgs))
	var thePkg *packages.Package
	for i, pkg := range apiPkgs {
		pkgPaths[i] = pkg.PkgPath
		if hasBase(pkg.PkgPath, pkgPath) {
			if thePkg != nil {
				err = generator.Errorf(nil, "ambiguous path (%v)", pkgPath)
				return
			}
			thePkg = pkg
		}
	}
	if thePkg == nil {
		err = generator.Errorf(nil, "invalid package path (%v not found in %v)", pkgPath, strings.Join(pkgPaths, ","))
		return
	}
	if !reName.MatchString(name) {
		err = generator.Errorf(nil, "invalid type name (%v)", name)
		return
	}
	return objName{pkg: thePkg.PkgPath, name: name}, opts, nil
}

func hasBase(pkgPath, tail string) bool {
	return pkgPath == tail ||
		strings.HasSuffix(pkgPath, tail) && pkgPath[len(pkgPath)-len(tail)-1] == '/'
}

func generateConverts(p generator.Printer, apiObjMap map[objName]*objMap) (count int, _ error) {
	list := make([]objName, 0, len(apiObjMap))
	for objName, obj := range apiObjMap {
		if len(obj.gens) != 0 {
			list = append(list, objName)
		}
	}
	// sort by package path then sort by name
	sort.Slice(list, func(i, j int) bool {
		if list[i].pkg < list[j].pkg {
			return true
		}
		if list[i].pkg > list[j].pkg {
			return false
		}
		return list[i].name < list[j].name
	})

	for _, objName := range list {
		m := apiObjMap[objName]
		if len(m.gens) == 0 {
			continue
		}
		w(p, "//-- convert %v.%v --//\n", m.src.Pkg().Path(), m.src.Name())
		for _, g := range m.gens {
			var err error
			switch g.mode {
			case ModeType:
				err = generateConvertType(p, m.src, g.obj)
			case ModeCreate:
				err = generateCreate(p, m.src, g.obj)
			case ModeUpdate:
				err = generateUpdate(p, m.src, g.obj, g.opts)
			default:
				panic("unexpected")
			}
			if err != nil {
				return count, generator.Errorf(err, "can not convert between %v and %v: %v", g.obj.Object.Name(), m.src.Object.Name(), err)
			}
			count++
		}
	}
	return count, nil
}

func generateConvertType(p generator.Printer, src, dst generator.Object) error {
	if err := generateConvertTypeImpl(p, src, dst); err != nil {
		return err
	}
	return generateConvertTypeImpl(p, dst, src)
}

func generateConvertTypeImpl(p generator.Printer, in, out generator.Object) error {
	inSt := validateStruct(in.Object)
	outSt := validateStruct(out.Object)
	inType := p.TypeString(in.Type())
	outType := p.TypeString(out.Type())
	fields := make([]fieldConvert, 0, outSt.NumFields())
	for i, n := 0, outSt.NumFields(); i < n; i++ {
		outField := outSt.Field(i)
		inField := matchField(outField, inSt)
		fields = append(fields, fieldConvert{
			Arg: inField,
			Out: outField,
		})
	}
	vars := map[string]interface{}{
		"InStr":   strings.ReplaceAll(inType, ".", "_"),
		"OutStr":  strings.ReplaceAll(outType, ".", "_"),
		"InType":  inType,
		"OutType": outType,
		"Fields":  fields,
	}
	return tplConvertType.Execute(p, vars)
}

func generateCreate(p generator.Printer, base, arg generator.Object) error {
	baseSt := validateStruct(base.Object)
	argSt := validateStruct(arg.Object)
	baseType := p.TypeString(base.Type())
	argType := p.TypeString(arg.Type())
	fields := make([]fieldConvert, 0, baseSt.NumFields())
	for i, n := 0, baseSt.NumFields(); i < n; i++ {
		baseField := baseSt.Field(i)
		argField := matchField(baseField, argSt)
		fields = append(fields, fieldConvert{
			Arg: argField,
			Out: baseField,
		})
	}
	vars := map[string]interface{}{
		"ArgStr":   strings.ReplaceAll(argType, ".", "_"),
		"ArgType":  argType,
		"BaseType": baseType,
		"Fields":   fields,
	}
	return tplCreate.Execute(p, vars)
}

func generateUpdate(p generator.Printer, base, arg generator.Object, opts options) error {
	baseSt := validateStruct(base.Object)
	argSt := validateStruct(arg.Object)
	baseType := p.TypeString(base.Type())
	argType := p.TypeString(arg.Type())
	fields := make([]fieldConvert, 0, baseSt.NumFields())
	for i, n := 0, baseSt.NumFields(); i < n; i++ {
		baseField := baseSt.Field(i)
		argField := matchField(baseField, argSt)
		fields = append(fields, fieldConvert{
			Arg: argField,
			Out: baseField,

			IsIdentifier: contains(opts.identifiers, baseField.Name()),
		})
	}
	vars := map[string]interface{}{
		"ArgStr":   strings.ReplaceAll(argType, ".", "_"),
		"ArgType":  argType,
		"BaseType": baseType,
		"Fields":   fields,
	}
	return tplUpdate.Execute(p, vars)
}

func contains(ss []string, item string) bool {
	for _, s := range ss {
		if s == item {
			return true
		}
	}
	return false
}

func matchField(baseField *types.Var, st *types.Struct) *types.Var {
	for i, n := 0, st.NumFields(); i < n; i++ {
		field := st.Field(i)
		if field.Name() == baseField.Name() {
			return field
		}
	}
	return nil
}

func w(w io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(w, format, args...)
	if err != nil {
		panic(err)
	}
}
