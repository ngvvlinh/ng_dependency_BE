package swagger

import (
	"strings"

	"github.com/go-openapi/jsonreference"
	"github.com/go-openapi/spec"
)

const prefixDef = "#/definitions/"

type Definitions map[string]spec.Schema

func (defs Definitions) ByURL(urlStr string) *spec.Schema {
	if !strings.HasPrefix(urlStr, prefixDef) {
		return nil
	}
	id := strings.TrimPrefix(urlStr, prefixDef)
	def, ok := defs[id]
	if !ok {
		return nil
	}
	return &def
}

func referenceByID(id string) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.Ref{Ref: jsonreference.MustCreateRef(prefixDef + id)},
		},
	}
}

// convert $ref and description to allOf
func convertSchema(schema spec.Schema, definitions Definitions) spec.Schema {
	u := schema.SchemaProps.Ref.GetURL()
	if u == nil {
		return schema
	}
	if schema.Description == "" {
		return schema
	}
	def := definitions.ByURL(u.String())

	description := schema.Description
	schema.Description = ""
	if !strings.HasSuffix(description, "\n") {
		description += "\n"
	}
	if def.Description != "" {
		description = description + "\n" + def.Description
	}

	extSchema := spec.Schema{
		SchemaProps: spec.SchemaProps{
			Description: description,
		},
	}
	result := spec.Schema{
		SchemaProps: spec.SchemaProps{
			AllOf: []spec.Schema{extSchema, schema},
		},
	}
	return result
}

func fillSchemaDesc(schema *spec.Schema, desc ItemDescription) {
	schema.Description = desc.FormattedDescription
}

func coalesce(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
