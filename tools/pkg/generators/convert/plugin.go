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
	// collect all converting packages
	var generatingPackages []*generatingPackage
	pkgs := ng.GeneratingPackages()
	for _, pkg := range pkgs {
		p, err := preparePackage(ng, pkg)
		if err != nil {
			return err
		}
		generatingPackages = append(generatingPackages, p)
	}

	// find all package pairs
	pkgPairs := make(map[pkgPair]*generatingPackage)
	for _, gpkg := range generatingPackages {
		for _, step := range gpkg.steps {
			for _, argPkg := range step.argPkgs {
				for _, outPkg := range step.outPkgs {
					pair0 := pkgPair{argPkg.PkgPath, outPkg.PkgPath}
					pair1 := pkgPair{outPkg.PkgPath, argPkg.PkgPath}
					if pkgPairs[pair0] != nil {
						return generator.Errorf(nil, "multiple packages with same conversion %v->%v (%v and %v)",
							argPkg.PkgPath, outPkg.PkgPath, pkgPairs[pair0].gpkg.PkgPath, gpkg.gpkg.PkgPath)
					}
					pkgPairs[pair0] = gpkg
					pkgPairs[pair1] = gpkg
				}
			}
		}
	}

	// find all custom conversion functions, they must be declared in the conversion packages
	convPairs = make(map[pair]*conversionFunc)
	for _, gpkg := range generatingPackages {
		for _, object := range gpkg.gpkg.Objects() {
			fn, ok := object.Object.(*types.Func)
			if !ok {
				continue
			}
			sign := fn.Type().(*types.Signature)
			if sign.Recv() != nil {
				continue
			}
			mode, arg, out, err := validateConvertFunc(fn)
			if err != nil {
				ll.V(2).Debugf("error in function %v.%v: err", fn.Pkg().Path(), fn.Name(), err)
				return err
			}
			if mode == 0 {
				ll.V(2).Debugf("ignore function %v.%v because it is not a recognized signature format", fn.Pkg().Path(), fn.Name())
				continue
			}
			pkgPair := pkgPair{
				ArgPkg: arg.Pkg().Path(),
				OutPkg: out.Pkg().Path(),
			}
			if gpkg1 := pkgPairs[pkgPair]; gpkg1 != nil && gpkg1 != gpkg {
				return generator.Errorf(nil,
					"function %v which converts from %v to %v must be defined in %v (found in %v)",
					fn.Name(), arg.Name(), out.Name(),
					gpkg.gpkg.PkgPath, gpkg1.gpkg.PkgPath)
			}
			pair, _, _ := getPairWithPointer(arg, out)
			if !pair.valid {
				ll.V(2).Debugf("ignore function %v.%v because its params are not pointer to named type", fn.Pkg().Path(), fn.Name())
				continue
			}
			if convPairs[pair] != nil {
				return generator.Errorf(nil,
					"duplicated conversion functions from %v to %v (function %v and %v)",
					arg.Type().String(), out.Type().String(), convPairs[pair].Func.Name(), fn.Name())
			}
			convPairs[pair] = &conversionFunc{
				pair: pair,
				Obj:  object,
				Func: fn,
				Mode: mode,
			}
		}
	}

	// generate
	for _, gpkg := range generatingPackages {
		currentPrinter = gpkg.gpkg.Generate()
		_, err := generateConverts(currentPrinter, convPairs, gpkg.objMap)
		if err != nil {
			return err
		}
	}
	return nil
}

type generatingPackage struct {
	gpkg   *generator.GeneratingPackage
	objMap map[objName]*objMap
	steps  []*generatingPackageStep
}

type generatingPackageStep struct {
	outPkgs []*packages.Package
	argPkgs []*packages.Package
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
	Arg *types.Var
	Out *types.Var

	IsIdentifier bool // for updating only
}

type pkgPair struct {
	ArgPkg string
	OutPkg string
}

type pair struct {
	valid bool
	Arg   objName
	Out   objName
}

type conversionFunc struct {
	pair
	Obj  generator.Object
	Func *types.Func
	Mode int
}

func (o objName) String() string {
	if o.pkg == "" {
		return o.name
	}
	return o.pkg + "." + o.name
}

func preparePackage(ng generator.Engine, gpkg *generator.GeneratingPackage) (*generatingPackage, error) {
	result := &generatingPackage{
		gpkg:   gpkg,
		objMap: make(map[objName]*objMap),
	}
	for _, d := range gpkg.Directives {
		if d.Cmd != Command {
			continue
		}
		apiPkgPaths, toPkgPaths, err := parseConvertDirective(d)
		if err != nil {
			return nil, err
		}
		step, err := generatePackageStep(ng, gpkg.Generate(), result.objMap, apiPkgPaths, toPkgPaths)
		if err != nil {
			return nil, err
		}
		result.steps = append(result.steps, step)
	}
	if len(result.steps) == 0 {
		return nil, generator.Errorf(nil, "invalid directive (must in format pkg1 -> pkg2)")
	}
	return result, nil
}

func generatePackageStep(ng generator.Engine, p generator.Printer, apiObjMap map[objName]*objMap, apiPkgPaths, toPkgPaths []string) (*generatingPackageStep, error) {
	ll.V(1).Debugf("convert from %v to %v", strings.Join(apiPkgPaths, ","), strings.Join(toPkgPaths, ","))

	var result generatingPackageStep
	flagSelf, err := validateEquality(apiPkgPaths, toPkgPaths)
	if err != nil {
		return nil, err
	}
	flagAuto := !flagSelf && len(apiPkgPaths) == 1 && len(toPkgPaths) == 1

	apiPkgs := make([]*packages.Package, len(apiPkgPaths))
	for i, pkgPath := range apiPkgPaths {
		apiPkgs[i] = ng.PackageByPath(pkgPath)
		if apiPkgs[i] == nil {
			return nil, generator.Errorf(nil, "can not find package %v", pkgPath)
		}
	}
	toPkgs := make([]*packages.Package, len(toPkgPaths))
	for i, pkgPath := range toPkgPaths {
		toPkgs[i] = ng.PackageByPath(pkgPath)
		if toPkgs[i] == nil {
			return nil, generator.Errorf(nil, "can not find package %v", pkgPath)
		}
	}

	for _, pkg := range apiPkgs {
		apiObjs := ng.ObjectsByPackage(pkg)
		for _, obj := range apiObjs {
			if !obj.Object.Exported() {
				continue
			}
			objName := objName{pkg: pkg.PkgPath, name: obj.Object.Name()}
			if apiObjMap[objName] == nil {
				apiObjMap[objName] = &objMap{src: obj}
			}
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
				return nil, err
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
				}

			} else if apiObjMap[name] == nil {
				return nil, generator.Errorf(nil, "type %v not found (directive %v)", name, raw)
			}
			if (name == objName{}) {
				continue
			}

			m := apiObjMap[name]
			if s := validateStruct(m.src.Object); s == nil {
				return nil, generator.Errorf(nil, "%v is not a struct", m.src.Object.Name())
			}
			m.gens = append(m.gens, objGen{mode: mode, obj: obj, opts: opts})
		}
	}
	return &result, nil
}

func validateConvertFunc(fn *types.Func) (mode int, arg, out *types.Var, err error) {
	sign := fn.Type().(*types.Signature)
	params, results := sign.Params(), sign.Results()
	switch {
	case params.Len() == 1 && results.Len() == 1:
		arg = params.At(0)
		out = results.At(0)
		return 1, arg, out, nil

	case params.Len() == 2 && results.Len() == 0:
		arg = params.At(0)
		out = params.At(1)
		return 2, arg, out, nil

	case params.Len() == 2 && results.Len() == 0:
		if params.At(1).Type() == results.At(0).Type() {
			arg = params.At(0)
			out = params.At(1)
			return 3, arg, out, nil
		}
	}
	// ignore unrecognized functions
	return 0, nil, nil, nil
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

func generateConverts(p generator.Printer, convPair map[pair]*conversionFunc, apiObjMap map[objName]*objMap) (count int, _ error) {
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
		"InStr":                strings.ReplaceAll(inType, ".", "_"),
		"OutStr":               strings.ReplaceAll(outType, ".", "_"),
		"InType":               inType,
		"OutType":              outType,
		"Fields":               fields,
		"CustomConversionMode": 0,
	}
	if conv := convPairs[getPair(in, out)]; conv != nil {
		vars["CustomConversionMode"] = conv.Mode
		funcType := conv.Obj.Name()
		alias := p.Qualifier(conv.Obj.Pkg())
		if alias != "" {
			funcType = alias + "." + funcType
		}
		vars["CustomConversionFuncType"] = funcType
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
	identCount := 0
	for i, n := 0, baseSt.NumFields(); i < n; i++ {
		baseField := baseSt.Field(i)
		argField := matchField(baseField, argSt)
		isIdentifier := contains(opts.identifiers, baseField.Name())
		if isIdentifier {
			identCount++
		}
		fields = append(fields, fieldConvert{
			Arg: argField,
			Out: baseField,

			IsIdentifier: isIdentifier,
		})
	}
	if identCount != len(opts.identifiers) {
		return fmt.Errorf("update %v: identifier not found (%v)", arg.Name(), strings.Join(opts.identifiers, ","))
	}
	vars := map[string]interface{}{
		"ArgStr":   strings.ReplaceAll(argType, ".", "_"),
		"ArgType":  argType,
		"BaseType": baseType,
		"Fields":   fields,
	}
	return tplUpdate.Execute(p, vars)
}

func validateCompatible(arg, out types.Object) bool {
	if arg.Type() == out.Type() {
		return true
	}
	slice0, ok0 := arg.Type().(*types.Slice)
	slice1, ok1 := out.Type().(*types.Slice)
	return ok0 && ok1 && slice0.Elem() == slice1.Elem()
}

func validateSliceToPointerNamed(obj types.Type) *types.Named {
	if typ, ok := obj.(*types.Named); ok {
		obj = typ.Underlying()
	}
	slice, ok := obj.(*types.Slice)
	if !ok {
		return nil
	}
	ptr, ok := slice.Elem().(*types.Pointer)
	if !ok {
		return nil
	}
	named, _ := ptr.Elem().(*types.Named)
	return named
}

func validatePointerToNamed(obj types.Type) *types.Named {
	typ, ok := obj.(*types.Pointer)
	if !ok {
		return nil
	}
	named, _ := typ.Elem().(*types.Named)
	return named
}

func getPairWithSlice(arg, out types.Object) (result pair, argNamed, outNamed *types.Named) {
	argNamed = validateSliceToPointerNamed(arg.Type())
	outNamed = validateSliceToPointerNamed(out.Type())
	if argNamed == nil {
		return
	}
	if outNamed == nil {
		return
	}
	result = pair{
		valid: true,
		Arg:   getObjName(argNamed),
		Out:   getObjName(outNamed),
	}
	return
}

func getPairWithPointer(arg, out types.Object) (result pair, argNamed, outNamed *types.Named) {
	argNamed = validatePointerToNamed(arg.Type())
	outNamed = validatePointerToNamed(out.Type())
	if argNamed == nil {
		ll.V(3).Debugf("ignore type %v because it is not a pointer to a named type", arg.Type())
		return
	}
	if outNamed == nil {
		ll.V(3).Debugf("ignore type %v because it is not a pointer to a named type", out.Type())
		return
	}
	result = pair{
		valid: true,
		Arg:   getObjName(argNamed),
		Out:   getObjName(outNamed),
	}
	return
}

func getPair(arg, out types.Object) pair {
	argNamed, _ := arg.Type().(*types.Named)
	outNamed, _ := out.Type().(*types.Named)
	if argNamed == nil || outNamed == nil {
		return pair{}
	}
	return pair{
		valid: true,
		Arg:   getObjName(argNamed),
		Out:   getObjName(outNamed),
	}
}

func getObjName(named *types.Named) objName {
	return objName{
		pkg:  named.Obj().Pkg().Path(),
		name: named.Obj().Name(),
	}
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
