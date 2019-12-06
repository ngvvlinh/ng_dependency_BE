package swagger

import (
	"encoding/json"
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
	"golang.org/x/tools/go/packages"

	"etop.vn/backend/tools/pkg/gen"
	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/generators/api/defs"
	"etop.vn/backend/tools/pkg/generators/api/parse"
	"etop.vn/backend/tools/pkg/generators/apix"
	"etop.vn/backend/tools/pkg/genutil"
	"etop.vn/common/l"
)

var ll = l.New()
var info *parse.Info

var _ generator.Plugin = &plugin{}

type plugin struct {
	generator.Filterer
}

type Opts struct {
	apix.Opts
	Description string
}

func New() generator.Plugin {
	return &plugin{
		Filterer: generator.FilterByCommand("gen:apix"),
	}
}

func (p *plugin) Name() string { return "swagger" }

func (p *plugin) Generate(ng generator.Engine) error {
	info = parse.NewInfo(ng)
	return ng.GenerateEachPackage(p.generatePackage)
}

func (p *plugin) generatePackage(ng generator.Engine, pkg *packages.Package, _ generator.Printer) (_err error) {
	pkgDirectives := ng.GetDirectivesByPackage(pkg)
	basePath := pkgDirectives.GetArg("gen:apix:base-path")
	if basePath == "" {
		basePath = "/api"
	}
	docPath := pkgDirectives.GetArg("gen:swagger:doc-path")
	if docPath == "" {
		return generator.Errorf(nil, "no doc-path for pkg %v", pkg.Name)
	}
	description, err := parseDescription(pkg, pkgDirectives)
	if err != nil {
		return err
	}
	opts := Opts{Description: description}
	opts.BasePath = basePath

	services, err := parse.Services(ng, pkg, []defs.Kind{defs.KindService})
	if err != nil {
		return err
	}
	swaggerDoc, err := GenerateSwagger(ng, opts, services)
	if err != nil {
		return generator.Errorf(err, "generate swagger: %v", err)
	}
	{
		dir := filepath.Join(gen.ProjectPath(), "doc", docPath)
		filename := filepath.Join(dir, "swagger.json")
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer func() {
			err := f.Close()
			if _err == nil {
				_err = err
			}
		}()
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(swaggerDoc); err != nil {
			return generator.Errorf(nil, "generate swagger: %v", err)
		}
	}
	return nil
}

func GenerateSwagger(ng generator.Engine, opts Opts, services []*defs.Service) (*spec.SwaggerProps, error) {
	definitions := map[string]spec.Schema{}
	pathItems := map[string]spec.PathItem{}
	for _, service := range services {
		for _, method := range service.Methods {
			sign := method.Method.Type().(*types.Signature)
			requestRef := getReference(ng, definitions, sign.Params().At(1).Type())
			responseRef := getReference(ng, definitions, sign.Results().At(0).Type())

			apiPath := opts.BasePath + service.APIPath + "/" + method.Name
			pathItem := spec.PathItem{
				Refable:          spec.Refable{},
				VendorExtensible: spec.VendorExtensible{},
				PathItemProps: spec.PathItemProps{
					Post: &spec.Operation{
						VendorExtensible: spec.VendorExtensible{},
						OperationProps: spec.OperationProps{
							Tags: []string{service.Name},
							ID:   method.Name,
							Parameters: []spec.Parameter{
								{
									Refable: spec.Refable{},
									ParamProps: spec.ParamProps{
										Name:     "body",
										In:       "body",
										Required: true,
										Schema: &spec.Schema{
											SchemaProps: spec.SchemaProps{
												Ref: requestRef,
											},
										},
									},
								},
							},
							Responses: &spec.Responses{
								ResponsesProps: spec.ResponsesProps{
									StatusCodeResponses: map[int]spec.Response{
										200: {
											ResponseProps: spec.ResponseProps{
												Description: "A successful response",
												Schema: &spec.Schema{
													SchemaProps: spec.SchemaProps{
														Ref: responseRef,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}
			pathItems[apiPath] = pathItem
		}
	}

	paths := &spec.Paths{Paths: pathItems}
	info := &spec.Info{
		InfoProps: spec.InfoProps{
			Version:     "v1",
			Title:       "etop API",
			Description: opts.Description,
		},
	}
	var tags []spec.Tag
	for _, s := range services {
		tags = append(tags, spec.Tag{
			TagProps: spec.TagProps{
				Name: s.Name,
			},
		})
	}
	sort.Slice(tags, func(i, j int) bool { return tags[i].Name < tags[j].Name })

	result := &spec.SwaggerProps{
		Swagger:     "2.0",
		Info:        info,
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json"},
		Schemes:     []string{"http", "https"},
		Paths:       paths,
		Definitions: definitions,
		Tags:        tags,
	}
	return result, nil
}

var mapComp2OrigPath = map[string]string{}
var mapOrig2CompPath = map[string]string{}

func getDefinitionID(typ types.Type) string {
	ptr, ok := typ.(*types.Pointer)
	if ok {
		typ = ptr.Elem()
	}
	named, ok := typ.(*types.Named)
	if !ok {
		panic(fmt.Sprintf("must be named type (got %v)", typ))
	}

	origPath := named.Obj().Pkg().Path()
	if compPath := mapOrig2CompPath[origPath]; compPath != "" {
		return compPath + named.Obj().Name()
	}

	compPath := origPath
	compPath = strings.TrimPrefix(compPath, "etop.vn/api/pb/")
	compPath = strings.TrimPrefix(compPath, "etop.vn/backend/")
	compPath = strings.TrimPrefix(compPath, "etop.vn/api/")
	compPath = strings.TrimPrefix(compPath, "etop.vn/")

	i := byte(10)
	compPath = rotatePath(compactPath(compPath), i)
	for mapComp2OrigPath[compPath] != "" {
		i += 7
		compPath = rotatePath(compPath, i)
	}
	mapComp2OrigPath[compPath] = origPath
	mapOrig2CompPath[origPath] = compPath
	return compPath + named.Obj().Name()
}

func compactPath(s string) string {
	var b [2]byte
	parts := strings.Split(s, "/")
	for _, p := range parts {
		b[0] += p[0]
		b[1] += p[len(p)-1]
	}
	return string(b[:])
}

func rotatePath(pkgPath string, n byte) string {
	p := []byte(pkgPath)
	for i := 0; i < len(p); i++ {
		p[i] = rotateChar(p[i], n)
		n += 3
	}
	return string(p)
}

func rotateChar(c byte, i byte) byte {
	c += i
	c = 'a' + (c-'a')%('z'-'a'+1)
	return c
}

func getReference(ng generator.Engine, definitions map[string]spec.Schema, typ types.Type) spec.Ref {
	typs, inner := genutil.ExtractType(typ)
	if typs[len(typs)-1] != genutil.Named {
		panic(fmt.Sprintf("must be named type (got %v)", inner))
	}
	id := getDefinitionID(inner)
	if _, ok := definitions[id]; !ok {
		parseSchema(ng, definitions, inner)
	}
	return spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/" + id)}
}

func parseSchema(ng generator.Engine, definitions map[string]spec.Schema, typ types.Type) spec.Schema {
	ll.V(3).Debugf("parse schema for type %v", typ)
	var inner types.Type
	switch {
	case info.IsTime(typ):
		return simpleType("string", "date-time")

	case info.IsSliceOfBytes(typ):
		return simpleType("string", "byte")

	case info.IsNullID(typ):
		return simpleType("string", "int64")

	case info.IsBasic(typ, &inner) || info.IsNullBasic(typ, &inner):
		switch inner.(*types.Basic).Kind() {
		case types.Bool:
			return simpleType("boolean", "")

		case types.Int:
			return simpleType("integer", "")

		case types.Int32:
			return simpleType("integer", "int32")

		case types.Int64:
			return simpleType("string", "int64")

		case types.Float32:
			return simpleType("integer", "float32")

		case types.Float64:
			return simpleType("integer", "float32")

		case types.String:
			return simpleType("string", "")
		}

	case info.IsNamedStruct(typ, &inner):
		id := getDefinitionID(typ)
		refSchema := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Ref: spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/" + id)},
			},
		}
		if _, ok := definitions[id]; ok {
			return refSchema
		}

		// placeholder to prevent infinite recursion
		definitions[id] = spec.Schema{}
		props := map[string]spec.Schema{}
		st := inner.(*types.Struct)

		// TODO: use message.Walk
		for i, n := 0, st.NumFields(); i < n; i++ {
			field := st.Field(i)
			if !field.Exported() {
				continue
			}
			jsonTag := parseJsonTag(st.Tag(i))
			switch jsonTag {
			case "":
				panic(fmt.Sprintf("no tag on field %v of struct %v", field.Name(), typ))
			case "-":
				continue
			}
			fieldSchema := parseSchema(ng, definitions, field.Type())
			props[jsonTag] = fieldSchema
		}
		s := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type:       spec.StringOrArray{"object"},
				Properties: props,
			},
		}
		definitions[id] = s
		return refSchema

	case info.IsArray(typ, &inner):
		refSchema := parseSchema(ng, definitions, inner)
		s := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"array"},
				Items: &spec.SchemaOrArray{
					Schema: &refSchema,
				},
			},
		}
		return s

	case info.IsMap(typ):
		m := typ.(*types.Map)
		elemSchema := parseSchema(ng, definitions, m.Elem())
		s := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				AdditionalProperties: &spec.SchemaOrBool{
					Schema: &elemSchema,
				},
			},
		}
		return s

	case info.IsEnum(typ):
		id := getDefinitionID(typ)
		refSchema := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Ref: spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/" + id)},
			},
		}
		if _, ok := definitions[id]; ok {
			return refSchema
		}

		enum := info.GetEnum(typ)
		var enumNames []interface{}
		for _, value := range enum.Values {
			enumNames = append(enumNames, enum.MapName[value])
		}

		var deprecatedEnumNames []string
		for name, value := range enum.MapValue {
			if enum.MapName[value] != name {
				deprecatedEnumNames = append(deprecatedEnumNames, name)
			}
		}

		s := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"string"},
				Enum: enumNames,
			},
		}
		if len(deprecatedEnumNames) != 0 {
			s.Description = fmt.Sprintf(`Deprecated values: "%v"`, strings.Join(deprecatedEnumNames, `", "`))
		}
		definitions[id] = s
		return refSchema

	case info.IsID(typ):
		return simpleType("string", "int64")

	case info.IsNamedInterface(typ, &inner):
		panic(fmt.Sprintf("oneof is not supported"))
	}

	panic(fmt.Sprintf("unsupported %v", typ))
}

func simpleType(typ, format string) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Type:   spec.StringOrArray{typ},
			Format: format,
		},
	}
}

func parseJsonTag(tag string) string {
	st, err := genutil.ParseStructTags(tag)
	if err != nil {
		panic(fmt.Sprintf("invalid tag %v", tag))
	}
	for _, t := range st {
		if t.Name == "json" {
			parts := strings.Split(t.Value, ",")
			return parts[0]
		}
	}
	return ""
}

func parseDescription(pkg *packages.Package, ds generator.Directives) (string, error) {
	filePath := ds.GetArg("gen:swagger:description")
	if filePath == "" {
		return "", nil
	}
	absPath := filepath.Join(filepath.Dir(pkg.GoFiles[0]), filePath)
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return "", generator.Errorf(err, "%v", err)
	}
	return strings.TrimSpace(string(data)), nil
}
