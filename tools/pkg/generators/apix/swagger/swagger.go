package swagger

import (
	"fmt"
	"go/types"
	"sort"
	"strings"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"

	"etop.vn/backend/tools/pkg/generator"
	"etop.vn/backend/tools/pkg/generators/apix/defs"
	. "etop.vn/backend/tools/pkg/genutil"
	"etop.vn/common/l"
)

const description = `API Documentation:
- [/doc/etop](/doc/etop) Shared API for managing login and account
- [/doc/shop](/doc/shop) API for shops
- [/doc/affiliate](/doc/affiliate) API for affiliates
- [/doc/services/affiliate](/doc/services/affiliate) API for affiliate service
- [/doc/integration](/doc/integration) API for shop to integrate with external partners
- [/doc/admin](/doc/admin) API for eTop admins
- [/doc/ext/shop](/doc/ext/shop) External API for shops
- [/doc/ext/partner](/doc/ext/partner) External API for partners
- [/doc/sadmin](/doc/sadmin) Special API for super admins
`

var ll = l.New()

func GenerateSwagger(ng generator.Engine, services []*defs.Service) (*spec.SwaggerProps, error) {
	initTypes(ng)
	definitions := map[string]spec.Schema{}
	pathItems := map[string]spec.PathItem{}
	for _, service := range services {
		for _, method := range service.Methods {
			sign := method.Method.Type().(*types.Signature)
			requestRef := getReference(ng, definitions, sign.Params().At(1).Type())
			responseRef := getReference(ng, definitions, sign.Results().At(0).Type())

			apiPath := service.BasePath + service.APIPath + "/" + method.Name
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
			Version: "v1",
			Title:   "etop API",
		},
	}
	// TODO: remove hard-code
	if services[0].BasePath == "/api" {
		info.Description = description
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
	typs, inner := ExtractType(typ)
	if typs[len(typs)-1] != Named {
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
	case isTime(typ):
		return simpleType("string", "date-time")

	case isSliceOfBytes(typ):
		return simpleType("string", "byte")

	case isNullID(typ):
		return simpleType("string", "int64")

	case isBasic(typ, &inner) || isNullBasic(typ, &inner):
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

	case isNamedStruct(typ, &inner):
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

	case isArray(typ, &inner):
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

	case isMap(typ):
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

	case isEnum(typ, &inner):
		id := getDefinitionID(typ)
		refSchema := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Ref: spec.Ref{Ref: jsonreference.MustCreateRef("#/definitions/" + id)},
			},
		}
		if _, ok := definitions[id]; ok {
			return refSchema
		}

		// read all enum
		named := skipPointer(typ).(*types.Named)
		var enumValues []interface{}
		objects := ng.GetObjectsByScope(named.Obj().Pkg().Scope())
		for _, obj := range objects {
			cnst, ok := obj.(*types.Const)
			if !ok {
				continue
			}
			if cnst.Type() != named {
				continue
			}
			parts := strings.SplitN(cnst.Name(), "_", 2)
			if len(parts) != 2 {
				panic(fmt.Sprintf("invalid enum constant %v", cnst.Name()))
			}
			val := parts[1]
			enumValues = append(enumValues, val)
		}

		s := spec.Schema{
			SchemaProps: spec.SchemaProps{
				Type: spec.StringOrArray{"string"},
				Enum: enumValues,
			},
		}
		definitions[id] = s
		return refSchema

	case isID(typ):
		return simpleType("string", "int64")

	case isNamedInterface(typ, &inner):
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
	st, err := ParseStructTags(tag)
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
