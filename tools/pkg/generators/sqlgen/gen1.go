package sqlgen

import (
	"go/types"
	"path/filepath"
	"reflect"
	"strings"

	"o.o/backend/tools/pkg/gen"
	"o.o/backend/tools/pkg/generator"
	"o.o/backend/tools/pkg/generators/sqlgen/filtergen"
)

type Option func(*genImpl, *typeDef)

func OptionSimple(st *types.Named) Option {
	return func(g *genImpl, _ *typeDef) {
		g.bases = append(g.bases, st)
		g.mapBase[st.String()] = true
	}
}

func OptionDerived(st, base *types.Named) Option {
	return func(_ *genImpl, def *typeDef) {
		def.base = base
	}
}

func OptionAs(alias string) Option {
	return func(_ *genImpl, def *typeDef) {
		def.all = false
		def.alias = alias
	}
}

func OptionJoin(baseType *types.Named, joinType *types.Named, joinKeyword string, joinAlias string, joinCond string) Option {
	return func(_ *genImpl, def *typeDef) {
		def.all = false
		def.joins = append(def.joins, &joinDef{
			JoinKeyword: joinKeyword,
			JoinAlias:   joinAlias,
			JoinCond:    joinCond,
			JoinType:    joinType,
			BaseType:    baseType,
		})
	}
}

func (g *genImpl) AddStruct(st *types.Named, opts ...Option) error {
	def, err := g.addStruct(st)
	if err != nil {
		return err
	}
	for _, opt := range opts {
		opt(g, def)
	}
	g.parseDef(g.Printer, def, st)
	return nil
}

func (g *genImpl) addStruct(st *types.Named) (*typeDef, error) {
	cols, excols, err := parseColumnsFromType(nil, st, unwrapNamedStruct(st))
	if err != nil {
		return nil, err
	}

	preloads := make([]*preloadDef, len(excols))
	for i, col := range excols {
		typ := col.fieldType
		desc := GetTypeDesc(g.Printer, typ)
		if !desc.Ptr && desc.Container == reflect.Slice &&
			desc.PtrElem && desc.Elem == reflect.Struct {
			// continue
		} else {
			return nil, generator.Errorf(nil, "Preload type must be slice of pointer to struct (got %v)", desc.TypeString)
		}

		if !strings.HasPrefix(desc.TypeString, "[]*") {
			return nil, generator.Errorf(nil, "Only support []* for preload type")
		}
		bareTypeStr := desc.TypeString[3:]

		preload := &preloadDef{
			TableName:     tableNameFromType(bareTypeStr),
			FieldType:     col.fieldType,
			FieldName:     col.FieldName,
			PluralTypeStr: plural(bareTypeStr),
			BaseType:      nil, // TODO
			Fkey:          col.fkey,
		}
		preloads[i] = preload
	}

	def := &typeDef{
		typ:      st,
		all:      true,
		cols:     cols,
		preloads: preloads,
		structs:  getStructsFromCols(cols),
	}
	for _, col := range cols {
		if col.timeLevel > def.timeLevel {
			def.timeLevel = col.timeLevel
			break
		}
	}
	return def, nil
}

func (g *genImpl) parseDef(p generator.Printer, def *typeDef, st *types.Named) {
	if def.base != nil {
		def.tableName = tableNameFromType(pr.TypeString(def.base))
	} else {
		def.tableName = tableNameFromType(pr.TypeString(st))
	}
	g.mapType[st.String()] = def

	// genfilters
	g.nAdd++
	if len(def.joins) == 0 {
		if g.genFilter == nil {
			var pkgName string
			var pkgPath string
			_ = types.TypeString(st, func(p *types.Package) string {
				pkgName = p.Name()
				pkgPath = p.Path()
				return p.Path()
			})

			// generate sqlstore/filter.gen.go
			if pkgName == "model" {
				if strings.HasPrefix(pkgPath, ".") {
					var err error
					pkgPath, err = filepath.Abs(pkgPath)
					if err != nil {
						panic(err)
					}
					pkgPath, err = filepath.Rel(gen.ProjectPath(), pkgPath)
					if err != nil {
						panic(err)
					}
					pkgPath = filepath.Join("o.o/backend", pkgPath)
				}
				g.genFilter = filtergen.NewGen(pkgPath)
			}
		}

		_cols := make([]*filtergen.ColumnDef, len(def.cols))
		for i, col := range def.cols {
			_cols[i] = &filtergen.ColumnDef{
				ColumnName: col.ColumnName,
				FieldName:  col.FieldName,
				FieldType:  col.fieldType,
				TypeDesc:   GetTypeDesc(p, col.fieldType),
			}
		}
		if g.genFilter != nil {
			g.genFilter.AddTable(def.tableName, pr.TypeString(st), _cols)
		}

	} else {
		substructs := make([]string, 0, len(def.joins)+1)
		substructs = append(substructs, pr.TypeString(def.base))
		for _, join := range def.joins {
			substructs = append(substructs, pr.TypeString(join.JoinType))
		}
		if g.genFilter != nil {
			g.genFilter.AddJoinTable(pr.TypeString(st), substructs)
		}
	}
}
