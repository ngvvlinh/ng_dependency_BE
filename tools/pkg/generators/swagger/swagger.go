package swagger

import (
	"encoding/json"
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
var currentInfo *parse.Info

var _ generator.Plugin = &plugin{}

type plugin struct {
}

type Opts struct {
	apix.Opts
	Description string
}

func New() generator.Plugin {
	return &plugin{}
}

func (p *plugin) Name() string { return "swagger" }

func (p *plugin) Filter(ng generator.FilterEngine) error {
	currentInfo = parse.NewInfo(ng)
	return generator.FilterByCommand("gen:apix").Filter(ng)
}

func (p *plugin) Generate(ng generator.Engine) error {
	currentInfo.Init(ng)
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
	description, err := parsePackageDescription(pkg, pkgDirectives)
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
			desc, err := parseItemDescription(ng, method.Method)
			if err != nil {
				return nil, err
			}
			sign := method.Method.Type().(*types.Signature)
			requestRef := getReference(ng, definitions, sign.Params().At(1).Type())
			responseRef := getReference(ng, definitions, sign.Results().At(0).Type())

			apiPath := opts.BasePath + service.APIPath + "/" + method.Name
			pathItem := spec.PathItem{
				PathItemProps: spec.PathItemProps{
					Post: &spec.Operation{
						OperationProps: spec.OperationProps{
							Description: desc.FormattedDescription,
							Deprecated:  desc.Deprecated,
							Summary:     method.Name,
							Tags:        []string{service.Name},
							ID:          getOperationID(method),
							Parameters: []spec.Parameter{
								{
									ParamProps: spec.ParamProps{
										Name:     "body",
										In:       "body",
										Required: true,
										Schema:   requestRef,
									},
								},
							},
							Responses: &spec.Responses{
								ResponsesProps: spec.ResponsesProps{
									StatusCodeResponses: map[int]spec.Response{
										200: {
											ResponseProps: spec.ResponseProps{
												Description: "A successful response",
												Schema:      responseRef,
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

func getCompactPath(origPath string) string {
	if compPath := mapOrig2CompPath[origPath]; compPath != "" {
		return compPath
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
	return compPath
}

func getDefinitionID(typ types.Type) string {
	typ = parse.UnwrapNull(parse.SkipPointer(typ))
	named, ok := typ.(*types.Named)
	if !ok {
		panic(fmt.Sprintf("must be named type (got %v)", typ))
	}

	pkgPath := named.Obj().Pkg().Path()
	return getCompactPath(pkgPath) + named.Obj().Name()
}

func getOperationID(m *defs.Method) string {
	origPath := m.Method.Pkg().Path() + "." + m.Service.FullName
	return getCompactPath(origPath) + m.Name
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

func getReference(ng generator.Engine, definitions map[string]spec.Schema, typ types.Type) *spec.Schema {
	typs, inner := genutil.ExtractType(typ)
	if typs[len(typs)-1] != genutil.Named {
		panic(fmt.Sprintf("must be named type (got %v)", inner))
	}
	id := getDefinitionID(inner)
	if _, ok := definitions[id]; !ok {
		obj := inner.(*types.Named).Obj()
		name := obj.Name()
		if obj.Pkg() != nil {
			name = obj.Pkg().Path() + "." + name
		}
		parseSchema(ng, name, definitions, inner)
	}
	return &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/" + id)},
		},
	}
}

func getReferenceByID(id string) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/" + id)},
		},
	}
}

func parseSchema(ng generator.Engine, path string, definitions map[string]spec.Schema, typ types.Type) spec.Schema {
	ll.V(3).Debugf("parse schema for %v (type %v)", path, typ)
	var inner types.Type
	switch {
	case currentInfo.IsTime(typ):
		return simpleType("string", "date-time")

	case currentInfo.IsSliceOfBytes(typ):
		return simpleType("string", "byte")

	case currentInfo.IsNullID(typ):
		return simpleType("string", "int64")

	case currentInfo.IsBasic(typ, &inner) || currentInfo.IsNullBasic(typ, &inner):
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

	case currentInfo.IsID(typ):
		return simpleType("string", "int64")

	case currentInfo.IsEnum(typ):
		enum := currentInfo.GetEnum(typ)
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

		id := getDefinitionID(typ)
		definitions[id] = s
		return getReferenceByID(id)

	case currentInfo.IsNamed(typ, &inner):
		id := getDefinitionID(typ)
		if _, ok := definitions[id]; ok {
			return getReferenceByID(id)
		}

		// types with custom +swagger directive
		schema, err := parseSchemaDirectives(ng, inner.(*types.Named))
		if err != nil {
			panic(err)
		}
		if schema != nil {
			definitions[id] = *schema
			return getReferenceByID(id)
		}

		switch {
		case currentInfo.IsNamedStruct(typ, &inner):
			// placeholder to prevent infinite recursion
			definitions[id] = spec.Schema{}
			props := map[string]spec.Schema{}
			st := inner.(*types.Struct)

			// TODO: use message.Walk
			var requiredFields []string
			for i, n := 0, st.NumFields(); i < n; i++ {
				field := st.Field(i)
				if !field.Exported() {
					continue
				}
				jsonTag := parseJsonTag(st.Tag(i))
				switch jsonTag {
				case "":
					panic(fmt.Sprintf("no json tag on field %v of struct %v", field.Name(), typ))
				case "-":
					continue
				}
				fieldSchema := parseSchema(ng, path+"."+field.Name(), definitions, field.Type())
				desc, err := parseItemDescription(ng, field)
				if err != nil {
					panic(fmt.Sprintf("parse comment on field %v of struct %v: %v", field.Name(), typ, err))
				}
				fieldSchema.Description = desc.FormattedDescription
				if desc.Required {
					requiredFields = append(requiredFields, jsonTag)
				}
				props[jsonTag] = fieldSchema
			}
			s := spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type:       spec.StringOrArray{"object"},
					Properties: props,
					Required:   requiredFields,
				},
			}
			definitions[id] = s
			return getReferenceByID(id)

		case currentInfo.IsNamedInterface(typ, &inner):
			panic(fmt.Sprintf("oneof is not supported"))
		}

	case currentInfo.IsArray(typ, &inner):
		refSchema := parseSchema(ng, path+"[]", definitions, inner)
		s := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"array"},
				Items: &spec.SchemaOrArray{
					Schema: &refSchema,
				},
			},
		}
		return s

	case currentInfo.IsMap(typ):
		m := typ.(*types.Map)
		elemSchema := parseSchema(ng, path+"[]", definitions, m.Elem())
		s := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"object"},
				AdditionalProperties: &spec.SchemaOrBool{
					Schema: &elemSchema,
				},
			},
		}
		return s

	}

	panic(fmt.Sprintf("unsupported %v (%v)", typ, path))
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

func parsePackageDescription(pkg *packages.Package, ds generator.Directives) (string, error) {
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

type ItemDescription struct {
	Summary              string
	Description          string
	FormattedDescription string
	Deprecated           bool
	DeprecatedText       string
	Required             bool
	Default              string
}

var reDeprecated = regexp.MustCompile(`(?i)((?:^|\n)@?(deprecated|required|default|todo):?)(?:([^\n]+))?\n`)

func parseItemDescription(ng generator.Engine, pos generator.Positioner) (ItemDescription, error) {
	cmt := ng.GetComment(pos)
	doc := cmt.Text()
	desc, err := parseItemDescriptionText(doc)
	if err != nil {
		return desc, generator.Errorf(err, "%v: %v", pos, err)
	}
	return desc, nil
}

func parseItemDescriptionText(doc string) (res ItemDescription, _ error) {
	res.Description = doc
	match := reDeprecated.FindAllStringSubmatch(doc, -1)
	if len(match) == 0 {
		res.FormattedDescription = res.Description
		return
	}
	formattedDescription := doc
	for _, parts := range match {
		if strings.Contains(parts[0], " ") && !strings.ContainsAny(parts[1], "@:") {
			return res, generator.Errorf(nil, "invalid keyword, must contain @ or : (%v)", strings.TrimSpace(parts[0]))
		}
		keyword := strings.ToLower(parts[2])
		switch keyword {
		case "deprecated":
			res.Deprecated = true
			res.DeprecatedText = strings.TrimSpace(parts[3])
			formattedDescription = strings.Replace(formattedDescription, parts[1], "\n**Deprecated:**", 1)

		case "required":
			res.Required = true
			formattedDescription = strings.Replace(formattedDescription, parts[1], "", 1)

		case "default":
			res.Default = strings.TrimSpace(parts[3])
			formattedDescription = strings.Replace(formattedDescription, parts[1], "\n**Default:**", 1)

		case "todo":
			formattedDescription = strings.Replace(formattedDescription, parts[1], "\n**TODO:**", 1)

		default:
			panic(fmt.Sprintf("unexpected (%v)", keyword))
		}
	}
	if res.Required && res.Default != "" {
		return res, generator.Errorf(nil, "required and default can not be used together")
	}
	res.FormattedDescription = strings.TrimLeft(formattedDescription, "\n")
	return
}

func parseSchemaDirectives(ng generator.Engine, typ *types.Named) (*spec.Schema, error) {
	desc, err := parseItemDescription(ng, typ.Obj())
	if err != nil {
		return nil, err
	}

	directives := ng.GetDirectives(typ.Obj())
	swaggerType := directives.GetArg("swagger:type")
	swaggerNullable := directives.GetArg("swagger:nullable")
	swaggerFormat := directives.GetArg("swagger:format")
	if swaggerType == "" {
		return nil, nil
	}

	var nullable bool
	if swaggerNullable != "" {
		nullable, err = strconv.ParseBool(swaggerNullable)
		if err != nil {
			return nil, generator.Errorf(nil, "type %v: invalid directive +swagger:nullable=%v", typ, swaggerNullable)
		}
	}

	return &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Description: desc.FormattedDescription,
			Type:        spec.StringOrArray{swaggerType},
			Nullable:    nullable,
			Format:      swaggerFormat,
		},
	}, nil
}
