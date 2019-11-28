package convert

import (
	"fmt"
	"go/types"
	"io"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

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
		Filterer:  generator.FilterByCommand(Command),
		Qualifier: genutil.Qualifier{},
	}
}

type plugin struct {
	generator.Filterer
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
		for _, obj := range gpkg.gpkg.GetObjects() {
			fn, ok := obj.(*types.Func)
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
				gpkg.ignoredFuncs = append(gpkg.ignoredFuncs, nameWithComment{
					Name:    fn.Name(),
					Comment: "not recognized",
				})
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
				gpkg.ignoredFuncs = append(gpkg.ignoredFuncs, nameWithComment{
					Name:    fn.Name(),
					Comment: "params are not pointer to named types",
				})
				ll.V(2).Debugf("ignore function %v.%v because its params are not pointer to named types", fn.Pkg().Path(), fn.Name())
				continue
			}
			if convPairs[pair] != nil {
				return generator.Errorf(nil,
					"duplicated conversion functions from %v to %v (function %v and %v)",
					arg.Type().String(), out.Type().String(), convPairs[pair].Func.Name(), fn.Name())
			}
			customConv := nameWithComment{
				Name:    fn.Name(),
				Comment: "not use, no conversions between params",
			}
			if convInUse(gpkg.objMap, pair) {
				customConv.Comment = "in use"
			}
			gpkg.customConvs = append(gpkg.customConvs, customConv)
			convPairs[pair] = &conversionFunc{
				pair: pair,
				Func: fn,
				Mode: mode,
			}
		}
	}

	// generate
	for _, gpkg := range generatingPackages {
		currentPrinter = gpkg.gpkg.GetPrinter()
		generateComments(currentPrinter, gpkg.customConvs, gpkg.ignoredFuncs)
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

	customConvs  []nameWithComment
	ignoredFuncs []nameWithComment
}

type generatingPackageStep struct {
	outPkgs []*packages.Package
	argPkgs []*packages.Package
}

type objMap struct {
	src  types.Object
	gens []objGen
}

type objGen struct {
	mode    string
	obj     types.Object
	opts    options
	convPkg *packages.Package
}

type nameWithComment struct {
	Name    string
	Comment string
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
	Mode int
	Func *types.Func

	ConverterPkg *packages.Package
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
	for _, d := range gpkg.GetDirectives() {
		if d.Cmd != Command {
			continue
		}
		apiPkgPaths, toPkgPaths, err := parseConvertDirective(d)
		if err != nil {
			return nil, err
		}
		step, err := generatePackageStep(ng, gpkg, gpkg.GetPrinter(), result.objMap, apiPkgPaths, toPkgPaths)
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

func generatePackageStep(ng generator.Engine, gpkg *generator.GeneratingPackage, p generator.Printer, apiObjMap map[objName]*objMap, apiPkgPaths, toPkgPaths []string) (*generatingPackageStep, error) {
	ll.V(1).Debugf("convert from %v to %v", strings.Join(apiPkgPaths, ","), strings.Join(toPkgPaths, ","))

	var result generatingPackageStep
	flagSelf, err := validateEquality(apiPkgPaths, toPkgPaths)
	if err != nil {
		return nil, err
	}
	flagAuto := !flagSelf && len(apiPkgPaths) == 1 && len(toPkgPaths) == 1

	apiPkgs := make([]*packages.Package, len(apiPkgPaths))
	for i, pkgPath := range apiPkgPaths {
		apiPkgs[i] = ng.GetPackageByPath(pkgPath)
		if apiPkgs[i] == nil {
			return nil, generator.Errorf(nil, "can not find package %v", pkgPath)
		}
	}
	toPkgs := make([]*packages.Package, len(toPkgPaths))
	for i, pkgPath := range toPkgPaths {
		toPkgs[i] = ng.GetPackageByPath(pkgPath)
		if toPkgs[i] == nil {
			return nil, generator.Errorf(nil, "can not find package %v", pkgPath)
		}
	}

	for _, pkg := range apiPkgs {
		apiObjs := ng.GetObjectsByPackage(pkg)
		for _, obj := range apiObjs {
			if !obj.Exported() {
				continue
			}
			objName := objName{pkg: pkg.PkgPath, name: obj.Name()}
			if apiObjMap[objName] == nil {
				apiObjMap[objName] = &objMap{src: obj}
			}
		}
	}

	if ll.Verbosed(3) {
		for objName, objMap := range apiObjMap {
			ll.V(3).Debugf("object %v: %v", objName, objMap.src.Type())
		}
	}

	for _, toPkg := range toPkgs {
		toObjs := ng.GetObjectsByPackage(toPkg)
		for _, obj := range toObjs {
			directives := ng.GetDirectives(obj)
			ll.V(2).Debugf("convert to object %v with directives %#v", obj.Name(), directives)
			if !obj.Exported() {
				continue
			}
			if s := validateStruct(obj); s == nil {
				continue
			}
			raw, mode, name, opts, err := parseWithMode(apiPkgs, directives)
			if err != nil {
				return nil, err
			}
			ll.V(3).Debugf("parsed type %v with mode %v", name, mode)
			if mode == "" {
				// automatically convert type with the same name
				if flagAuto {
					name = objName{apiPkgs[0].PkgPath, obj.Name()}
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
			if s := validateStruct(m.src); s == nil {
				return nil, generator.Errorf(nil, "%v is not a struct", m.src.Name())
			}
			m.gens = append(m.gens, objGen{
				mode:    mode,
				obj:     obj,
				opts:    opts,
				convPkg: gpkg.Package,
			})
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

	case params.Len() == 2 && results.Len() == 1:
		if validatePtrTypeEquality(params.At(1).Type(), results.At(0).Type()) {
			arg = params.At(0)
			out = params.At(1)
			return 3, arg, out, nil
		}
	}

	// ignore unrecognized functions
	return 0, nil, nil, nil
}

func validatePtrTypeEquality(t0, t1 types.Type) bool {
	ptr0, ok0 := t0.(*types.Pointer)
	ptr1, ok1 := t1.(*types.Pointer)
	if !ok0 || !ok1 {
		return false
	}
	return ptr0.Elem() == ptr1.Elem()
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
	if !strings.Contains(directive.Arg, "->") {
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

func generateComments(
	p generator.Printer,
	customConversions, ignoredFuncs []nameWithComment,
) {
	sort.Slice(customConversions, func(i, j int) bool {
		return customConversions[i].Name < customConversions[j].Name
	})
	sort.Slice(ignoredFuncs, func(i, j int) bool {
		return ignoredFuncs[i].Name < ignoredFuncs[j].Name
	})
	tp := tabwriter.NewWriter(p, 16, 4, 0, ' ', 0)
	w(p, "/*\n")
	w(p, "Custom conversions:")
	if len(customConversions) == 0 {
		w(p, " (none)\n")
	} else {
		w(p, "\n")
	}
	for _, c := range customConversions {
		w(tp, "    %v", c.Name)
		if c.Comment != "" {
			w(tp, "\t    // %v", c.Comment)
		}
		w(tp, "\n")
	}
	_ = tp.Flush()
	w(p, "\nIgnored functions:")
	if len(ignoredFuncs) == 0 {
		w(p, " (none)\n")
	} else {
		w(p, "\n")
	}
	for _, c := range ignoredFuncs {
		w(tp, "    %v\t    // %v\n", c.Name, c.Comment)
	}
	_ = tp.Flush()
	w(p, "*/\n")
}

func convInUse(apiObjMap map[objName]*objMap, pair pair) bool {
	for _, item := range [][2]objName{
		{pair.Arg, pair.Out},
		{pair.Out, pair.Arg},
	} {
		m := apiObjMap[item[0]]
		if m == nil {
			continue
		}
		for _, g := range m.gens {
			name := objName{
				pkg:  g.obj.Pkg().Path(),
				name: g.obj.Name(),
			}
			if name == item[1] {
				return true
			}
		}
	}
	return false
}

func generateConverts(
	p generator.Printer,
	convPair map[pair]*conversionFunc,
	apiObjMap map[objName]*objMap,
) (count int, _ error) {
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

	// populate convPair with auto conversions
	for _, objName := range list {
		m := apiObjMap[objName]
		for _, g := range m.gens {
			if g.mode != ModeType {
				continue
			}
			pairs := []pair{getPair(m.src, g.obj), getPair(g.obj, m.src)}
			for _, pair := range pairs {
				if convPairs[pair] == nil {
					convPairs[pair] = &conversionFunc{pair: pair}
				}
				convPairs[pair].ConverterPkg = g.convPkg
			}
		}
	}

	var conversions []map[string]interface{}
	for _, objName := range list {
		m := apiObjMap[objName]
		for _, g := range m.gens {
			{
				arg, out := g.obj, m.src
				conversion := map[string]interface{}{}
				includeBaseConversion(p, conversion, g.mode, arg, out)
				conversions = append(conversions, conversion)
			}
			if g.mode == ModeType {
				arg, out := m.src, g.obj
				conversion := map[string]interface{}{}
				includeBaseConversion(p, conversion, g.mode, arg, out)
				conversions = append(conversions, conversion)
			}
		}
	}
	{
		p.Import("conversion", "etop.vn/backend/pkg/common/conversion")
		vars := map[string]interface{}{
			"Conversions": conversions,
		}
		if err := tplRegister.Execute(p, vars); err != nil {
			return 0, err
		}
	}

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
				err = generateConvertType(p, g.obj, m.src)
			case ModeCreate:
				err = generateCreate(p, g.obj, m.src)
			case ModeUpdate:
				err = generateUpdate(p, g.obj, m.src, g.opts)
			default:
				panic("unexpected")
			}
			if err != nil {
				return count, generator.Errorf(err, "can not convert between %v and %v: %v", g.obj.Name(), m.src.Name(), err)
			}
			count++
		}
	}
	return count, nil
}

func generateConvertType(p generator.Printer, src, dst types.Object) error {
	if err := generateConvertTypeImpl(p, src, dst); err != nil {
		return err
	}
	return generateConvertTypeImpl(p, dst, src)
}

func generateConvertTypeImpl(p generator.Printer, in types.Object, out types.Object) error {
	inSt := validateStruct(in)
	outSt := validateStruct(out)
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
		"Fields": fields,
	}
	includeBaseConversion(p, vars, ModeType, in, out)
	includeCustomConversion(p, vars, in, out)
	return tplConvertType.Execute(p, vars)
}

func generateCreate(p generator.Printer, arg types.Object, out types.Object) error {
	outSt := validateStruct(out)
	argSt := validateStruct(arg)
	fields := make([]fieldConvert, 0, outSt.NumFields())
	for i, n := 0, outSt.NumFields(); i < n; i++ {
		outField := outSt.Field(i)
		argField := matchField(outField, argSt)
		fields = append(fields, fieldConvert{
			Arg: argField,
			Out: outField,
		})
	}
	vars := map[string]interface{}{
		"Fields": fields,
	}
	includeBaseConversion(p, vars, ModeCreate, arg, out)
	includeCustomConversion(p, vars, arg, out)
	return tplCreate.Execute(p, vars)
}

func generateUpdate(p generator.Printer, arg types.Object, out types.Object, opts options) error {
	outSt := validateStruct(out)
	argSt := validateStruct(arg)
	fields := make([]fieldConvert, 0, outSt.NumFields())
	identCount := 0
	for i, n := 0, outSt.NumFields(); i < n; i++ {
		outField := outSt.Field(i)
		argField := matchField(outField, argSt)
		isIdentifier := contains(opts.identifiers, outField.Name())
		if isIdentifier {
			identCount++
		}
		fields = append(fields, fieldConvert{
			Arg: argField,
			Out: outField,

			IsIdentifier: isIdentifier,
		})
	}
	if identCount != len(opts.identifiers) {
		return fmt.Errorf("update %v: identifier not found (%v)", arg.Name(), strings.Join(opts.identifiers, ","))
	}
	vars := map[string]interface{}{
		"Fields": fields,
	}
	includeBaseConversion(p, vars, ModeUpdate, arg, out)
	includeCustomConversion(p, vars, arg, out)
	return tplUpdate.Execute(p, vars)
}

func includeBaseConversion(p generator.Printer, vars map[string]interface{}, mode string, arg types.Object, out types.Object) {
	outType := p.TypeString(out.Type())
	argType := p.TypeString(arg.Type())
	vars["ArgStr"] = strings.ReplaceAll(argType, ".", "_")
	vars["OutStr"] = strings.ReplaceAll(outType, ".", "_")
	vars["ArgType"] = argType
	vars["OutType"] = outType

	switch mode {
	case ModeType:
		vars["Action"] = "Convert"
		vars["action"] = "convert"
	case ModeCreate, ModeUpdate:
		vars["Action"] = "Apply"
		vars["action"] = "apply"
	default:
		panic("unexpected")
	}
}

func includeCustomConversion(p generator.Printer, vars map[string]interface{}, arg types.Object, out types.Object) {
	vars["CustomConversionMode"] = 0
	if conv := convPairs[getPair(arg, out)]; conv != nil && conv.Func != nil {
		vars["CustomConversionMode"] = conv.Mode
		funcName := conv.Func.Name()
		alias := p.Qualifier(conv.Func.Pkg())
		if alias != "" {
			funcName = alias + "." + funcName
		}
		vars["CustomConversionFuncName"] = funcName
	}
}

func validateCompatible(arg, out types.Object) bool {
	if arg.Type() == out.Type() {
		return true
	}
	{
		// *Type
		ptr0, ok0 := arg.Type().(*types.Pointer)
		ptr1, ok1 := out.Type().(*types.Pointer)
		if ok0 && ok1 && ptr0.Elem() == ptr1.Elem() {
			ll.V(1).Debugf("*Type %v %v", arg, out)
			return true
		}
	}
	{
		// []Type
		slice0, ok0 := arg.Type().(*types.Slice)
		slice1, ok1 := out.Type().(*types.Slice)
		if ok0 && ok1 && slice0.Elem() == slice1.Elem() {
			ll.V(1).Debugf("[]Type %v %v", arg, out)
			return true
		}

		// []*Type
		if ok0 && ok1 {
			ptr0, ptrok0 := slice0.Elem().(*types.Pointer)
			ptr1, ptrok1 := slice1.Elem().(*types.Pointer)
			if ptrok0 && ptrok1 && ptr0.Elem() == ptr1.Elem() {
				ll.V(1).Debugf("[]*Type %v %v", arg, out)
				return true
			}
		}
	}
	return false
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
