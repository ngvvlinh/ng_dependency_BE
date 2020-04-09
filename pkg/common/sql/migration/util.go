package migration

import (
	"fmt"
	"strings"

	"etop.vn/backend/pkg/common/sql/cmsql"
)

type ColumnDef struct {
	ColumnName       string
	ColumnType       string
	ColumnDBType     string
	ColumnTag        string
	ColumnEnumValues []string
}

var mDBTypeAndModelType = map[string]string{
	"integer":                  "int4",
	"timestamp with time zone": "timestamptz",
	"bigint":                   "int8",
	"smallint":                 "int2",
	"character varying":        "text",
}

func GetColumnNamesAndTypes(db *cmsql.Database, tableName string) (map[string]string, error) {
	rows, err := db.
		Select("column_name, data_type").
		From("information_schema.columns").
		Where(fmt.Sprintf("table_name = '%s'", tableName)).
		Query()
	if err != nil {
		return nil, err
	}

	mapColumnNameAndType := make(map[string]string)
	var columnName, dataType string
	for rows.Next() {
		err := rows.Scan(&columnName, &dataType)
		if err != nil {
			return nil, err
		}
		if val, ok := mDBTypeAndModelType[dataType]; ok {
			dataType = val
		}
		mapColumnNameAndType[columnName] = dataType
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mapColumnNameAndType, nil
}

func Compare(db *cmsql.Database, tableName string, mapColumnDef map[string]ColumnDef, mapDB map[string]string) error {
	mapModel := convert(mapColumnDef)
	for col, typ := range mapModel {
		DBTyp, ok := mapDB[col]
		// Handle missing field
		// case enum:
		// 	+ create enum_type
		//  + add column
		// default:
		//  + add column
		if !ok {
			enumType := getEnumType(mapColumnDef[col])
			fmt.Printf("Column (DB) %q is Missing\n", col)
			if mapColumnDef[col].ColumnDBType == "enum" {
				if _, err := db.Exec(fmt.Sprintf(`
					DO $$
						BEGIN
							IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN
								create type %q as enum (%s);
							END IF;
						END
					$$;`, enumType, enumType, convertArrayStringsToString(mapColumnDef[col].ColumnEnumValues))); err != nil {
					return err
				}
				if _, err := db.Exec(fmt.Sprintf("ALTER TABLE %q ADD COLUMN IF NOT EXISTS \"%s\" %s;\n", tableName, col, enumType)); err != nil {
					return err
				}
			} else {
				if _, err := db.Exec(fmt.Sprintf("ALTER TABLE %q ADD COLUMN IF NOT EXISTS \"%s\" %s;\n", tableName, col, typ)); err != nil {
					return err
				}
			}
			continue
		}
		// Handle not the same type
		// case "ARRAY"
		if (DBTyp == "ARRAY" && !strings.HasSuffix(typ, "[]")) || (DBTyp != "ARRAY" && DBTyp != typ) {
			fmt.Printf("Column %q is not of the same type %q and %q\n", col, typ, DBTyp)
		}
	}

	// Handle Field (model) is missing
	for col, _ := range mapDB {
		if _, ok := mapModel[col]; !ok {
			fmt.Printf("Field (model) %s is missing\n", col)
		}
	}
	return nil
}

var mColumnDBType = map[string]string{
	"[]*struct": "jsonb",
	"[]struct":  "jsonb",
	"*struct":   "jsonb",
	"[]string":  "text[]",
	"string":    "text",
	"int64":     "int8",
	"[]int64":   "int8[]",
	"[]byte":    "jsonb",
	"enum":      "USER-DEFINED",
}

var mModelTypeAndDBType = map[string]string{
	"string":         "text",
	"[]string":       "text[]",
	"time.Time":      "timestamptz",
	"status3.Status": "int2",
	"status4.Status": "int2",
	"status5.Status": "int2",
	"bool":           "boolean",
	"int":            "int4",
	"int8":           "int2",
	"[]int":          "int8[]",
	"dot.NullBool":   "boolean",
}

func convert(mapColumnDef map[string]ColumnDef) map[string]string {
	result := make(map[string]string)
	for _, col := range mapColumnDef {
		columnDBType := ""
		if col.ColumnTag != "" {
			columnDBType = col.ColumnTag
		} else if val, ok := mColumnDBType[col.ColumnDBType]; ok {
			columnDBType = val
		} else if val, ok := mModelTypeAndDBType[col.ColumnType]; ok {
			columnDBType = val
		}
		result[col.ColumnName] = columnDBType
	}
	return result
}

func convertArrayStringsToString(args []string) (result string) {
	if len(args) == 0 {
		return ""
	}
	for _, val := range args {
		result += fmt.Sprintf("'%s',", val)
	}

	return result[:len(result)-1]
}

func getEnumType(columDef ColumnDef) (enumType string) {
	tag := columDef.ColumnTag
	if tag != "" {
		if strings.HasPrefix(tag, "enum") {
			tag = strings.TrimPrefix(tag, "enum(")
			tag = tag[:len(tag)-1]
		}

		return tag
	}
	return strings.Split(columDef.ColumnType, ".")[0]
}
