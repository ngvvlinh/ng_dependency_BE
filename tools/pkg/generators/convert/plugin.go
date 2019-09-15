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
	"etop.vn/common/l"
)

var ll = l.New()

const Command = "gen:convert"
const ModeType = "convert:type"
const ModeApply = "convert:apply"

func New() generator.Plugin {
	return &plugin{
		Filter: generator.FilterByCommand(Command),
	}
}

type plugin struct {
	generator.Filter
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
}

type objName struct {
	pkg  string
	name string
}

type fieldConvert struct {
	OutField *types.Var
	InField  *types.Var
}

func (o objName) String() string {
	if o.pkg == "" {
		return o.name
	}
	return o.pkg + "." + o.name
}

func generatePackage(ng generator.Engine, gpkg *generator.GeneratingPackage) error {
	flagInfo := false
	apiPkgPaths, toPkgPath, err := parseConvertDirectives(gpkg.Directives)
	if err != nil {
		return err
	}

	ll.V(1).Debugf("convert from %v to %v", strings.Join(apiPkgPaths, ","), toPkgPath)

	apiPkgs := make([]*packages.Package, len(apiPkgPaths))
	for i, pkgPath := range apiPkgPaths {
		apiPkgs[i] = ng.PackageByPath(pkgPath)
		if apiPkgs[i] == nil {
			return fmt.Errorf("can not find package %v", pkgPath)
		}
	}

	toPkg := ng.PackageByPath(toPkgPath)
	if toPkg == nil {
		return fmt.Errorf("can not find package %v", toPkgPath)
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
	toObjs := ng.ObjectsByPackage(toPkg)
	for _, obj := range toObjs {
		ll.V(2).Debugf("convert to object %v with directives %#v", obj.Object.Name(), obj.Directives)
		if !obj.Object.Exported() {
			continue
		}
		if s := validateStruct(obj.Object); s == nil {
			continue
		}
		mode, name, err := parseWithMode(apiPkgs, obj.Directives)
		if err != nil {
			return err
		}
		if mode == "" && len(apiPkgs) == 1 {
			name = objName{apiPkgs[0].PkgPath, obj.Object.Name()}
			if apiObjMap[name] == nil {
				continue
			}
			mode = ModeType

		} else if mode == "" {
			flagInfo = true
			continue

		} else if apiObjMap[name] == nil {
			return fmt.Errorf("type %v not found", name)
		}

		m := apiObjMap[name]
		if s := validateStruct(m.src.Object); s == nil {
			return fmt.Errorf("%v is not a struct", m.src.Object.Name())
		}
		m.gens = append(m.gens, objGen{mode: mode, obj: obj})
	}

	currentPrinter = gpkg.Generate()
	count, err := generateConverts(currentPrinter, apiObjMap)
	if err != nil {
		return err
	}
	if flagInfo && count == 0 {
		ll.Warn("no types generated (for multiple package conversion, must use convert:type to define mapping)")
	}
	return nil
}

func validateStruct(obj types.Object) *types.Struct {
	s, ok := obj.(*types.TypeName)
	if !ok {
		return nil
	}
	st, _ := s.Type().Underlying().(*types.Struct)
	return st
}

func parseConvertDirectives(directives []generator.Directive) (apiPkgs []string, toPkg string, err error) {
	for _, d := range directives {
		if d.Cmd != Command {
			continue
		}

		// parse "pkg1 -> pkg2"
		parts := strings.Split(d.Arg, "->")
		if len(parts) != 2 {
			err = fmt.Errorf("invalid directive (must in format pkg1 -> pkg2)")
			return
		}

		// parse "pkg1,pkg2"
		apiPkgs = strings.Split(parts[0], ",")
		for i := range apiPkgs {
			apiPkgs[i] = strings.TrimSpace(apiPkgs[i])
			if apiPkgs[i] == "" {
				err = fmt.Errorf("invalid directive (must in format pkg1 -> pkg2)")
				return
			}
		}

		// validate toPkg
		toPkg = strings.TrimSpace(parts[1])
		if toPkg == "" {
			err = fmt.Errorf("invalid directive (must in format pkg1 -> pkg2)")
			return
		}
		return apiPkgs, toPkg, nil
	}

	err = fmt.Errorf("invalid directive (must in format pkg1 -> pkg2)")
	return
}

var reName = regexp.MustCompile(`[A-Z][A-z0-9_]*`)

func parseWithMode(apiPkgs []*packages.Package, directives []generator.Directive) (mode string, _ objName, _ error) {
	for _, d := range directives {
		switch d.Cmd {
		case ModeApply:
			objName, err := parseTypeName(apiPkgs, d.Arg)
			return d.Cmd, objName, err

		case ModeType:
			objName, err := parseTypeName(apiPkgs, d.Arg)
			return d.Cmd, objName, err
		}
	}
	return
}

func parseTypeName(apiPkgs []*packages.Package, input string) (_ objName, err error) {
	var pkgPath, name string
	idx := strings.LastIndexByte(input, '.')
	if idx < 0 {
		pkgPath, name = "", input
		if !reName.MatchString(name) {
			err = fmt.Errorf("invalid type name (%v)", name)
			return
		}
		if len(apiPkgs) == 1 {
			return objName{pkg: pkgPath, name: name}, nil
		}
		err = fmt.Errorf("must provide path for multiple input packages (%v)", input)
		return
	}

	pkgPath, name = input[:idx], input[idx+1:]
	var thePkg *packages.Package
	for _, pkg := range apiPkgs {
		if hasBase(pkg.PkgPath, pkgPath) {
			if thePkg != nil {
				err = fmt.Errorf("ambiguous path (%v)", pkgPath)
				return
			}
			thePkg = pkg
		}
	}
	if thePkg == nil {
		err = fmt.Errorf("invalid package path (%v)", pkgPath)
		return
	}
	if !reName.MatchString(name) {
		err = fmt.Errorf("invalid type name (%v)", name)
		return
	}
	return objName{pkg: thePkg.PkgPath, name: name}, nil
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
		for _, g := range m.gens {
			var err error
			switch g.mode {
			case ModeType:
				err = generateConvertType(p, m.src, g.obj)
			case ModeApply:
				err = generateConvertApply(p, m.src, g.obj)
			default:
				panic("unexpected")
			}
			if err != nil {
				return count, fmt.Errorf("can not convert between %v and %v: %v", g.obj.Object.Name(), m.src.Object.Name(), err)
			}
			count++
		}
	}
	return count, nil
}

func generateConvertType(p generator.Printer, src, dst generator.Object) error {
	srcStruct := validateStruct(src.Object)
	dstStruct := validateStruct(dst.Object)
	inType := p.TypeString(src.Object.Type())
	outType := p.TypeString(dst.Object.Type())
	fields := make([]fieldConvert, 0, dstStruct.NumFields())
	for i, n := 0, dstStruct.NumFields(); i < n; i++ {
		dstField := dstStruct.Field(i)
		var srcField *types.Var
		for j, m := 0, srcStruct.NumFields(); j < m; j++ {
			field := srcStruct.Field(j)
			if field.Name() == dstField.Name() {
				srcField = field
				break
			}
		}
		fields = append(fields, fieldConvert{OutField: dstField, InField: srcField})
	}
	vars := map[string]interface{}{
		"InStr":   strings.ReplaceAll(inType, ".", ""),
		"OutStr":  strings.ReplaceAll(outType, ".", ""),
		"InType":  inType,
		"OutType": outType,
		"Fields":  fields,
	}
	return tplConvertType.Execute(p, vars)
}

func generateConvertApply(p generator.Printer, src, dst generator.Object) error {
	return nil
}

func w(w io.Writer, format string, args ...interface{}) {
	_, err := fmt.Fprintf(w, format, args...)
	if err != nil {
		panic(err)
	}
}
