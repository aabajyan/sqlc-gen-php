package core

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/sqlc-dev/plugin-sdk-go/metadata"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

type Params struct {
	ModelClass *ModelClass
}

func (v Params) isEmpty() bool {
	return len(v.ModelClass.Fields) == 0
}

func splitFieldsByDefault(fields []Field, withDefault func(Field) string, withoutDefault func(Field) string) (nonDefaults, defaults []string) {
	for _, f := range fields {
		if f.Default != "" {
			defaults = append(defaults, withDefault(f))
			continue
		}

		nonDefaults = append(nonDefaults, withoutDefault(f))
	}
	return
}

func withoutDefaultCallback(f Field) string {
	return f.Type.String() + " $" + f.Name
}

func withDefaultCallback(f Field) string {
	return withoutDefaultCallback(f) + " = " + f.Default
}

func (v Params) Args() string {
	if v.isEmpty() {
		return ""
	}

	fields := v.ModelClass.Fields
	nonDefaults, defaults := splitFieldsByDefault(
		fields,
		withoutDefaultCallback,
		withoutDefaultCallback,
	)

	out := append(nonDefaults, defaults...)
	return strings.Join(out, ", ")
}

func (v Params) ArgsWithDefaults() string {
	if v.isEmpty() {
		return ""
	}

	fields := v.ModelClass.Fields
	nonDefaults, defaults := splitFieldsByDefault(
		fields,
		withDefaultCallback,
		withoutDefaultCallback,
	)

	out := append(nonDefaults, defaults...)
	return strings.Join(out, ", ")
}

func (v Params) Bindings() string {
	if v.isEmpty() {
		return ""
	}

	var out []string
	for _, f := range v.ModelClass.Fields {
		if f.Type.IsJSON() {
			out = append(out, fmt.Sprintf("json_encode($%s)", f.Name))
			continue
		}

		if f.Type.IsBoolean() {
			out = append(out, fmt.Sprintf("($%s ? 1 : 0)", f.Name))
			continue
		}

		out = append(out, fmt.Sprintf("$%s", f.Name))
	}

	return "[" + strings.Join(out, ", ") + "]"
}

func pdoRowMapping(t phpType, idx int) string {
	if t.IsJSON() {
		return fmt.Sprintf(`json_decode($row[%d], true) ?? []`, idx)
	}

	return fmt.Sprintf(`$row[%d]`, idx)
}

func (v QueryValue) PDOFetchMode() string {
	if !v.IsClass() {
		return "\\PDO::FETCH_COLUMN"
	}

	return "\\PDO::FETCH_NUM"
}

func (v QueryValue) ResultSet() string {
	if !v.IsClass() {
		return "$row"
	}

	var out []string
	for idx, f := range v.Struct.Fields {
		out = append(out, pdoRowMapping(f.Type, idx))
	}

	ret := strings.Join(out, ", ")
	return ret
}

func dataClassName(name string) string {
	out := ""
	for _, p := range strings.Split(name, "_") {
		out += cases.Title(language.English).String(p)
	}

	return out
}

func memberName(name string) string {
	return sdk.LowerTitle(dataClassName(name))
}

func BuildDataClasses(req *plugin.GenerateRequest) []*ModelClass {
	var structs []*ModelClass
	for _, schema := range req.Catalog.Schemas {
		if schema.Name == "pg_catalog" || schema.Name == "information_schema" {
			continue
		}

		for _, table := range schema.Tables {
			var tableName string
			if schema.Name == req.Catalog.DefaultSchema {
				tableName = table.Rel.Name
			} else {
				tableName = schema.Name + "_" + table.Rel.Name
			}

			structName := dataClassName(tableName)
			s := ModelClass{
				Table:   plugin.Identifier{Schema: schema.Name, Name: table.Rel.Name},
				Name:    structName,
				Comment: table.Comment,
			}

			for _, column := range table.Columns {
				s.Fields = append(s.Fields, Field{
					OriginalColumnName: column.Name,
					Name:               memberName(column.Name),
					Type:               makePhpTypeFromSqlcColumn(req, column),
					Comment:            column.Comment,
				})
			}
			structs = append(structs, &s)
		}
	}

	if len(structs) > 0 {
		sort.Slice(structs, func(i, j int) bool { return structs[i].Name < structs[j].Name })
	}
	return structs
}

func makePhpTypeFromSqlcColumn(req *plugin.GenerateRequest, col *plugin.Column) phpType {
	typ := mapSqlColumnTypeToPhpType(req, col)
	return phpType{
		Name:     typ,
		IsArray:  col.IsSqlcSlice,
		IsNull:   !col.NotNull,
		DataType: sdk.DataType(col.Type),
		Engine:   req.Settings.Engine,
	}
}

func mapSqlColumnTypeToPhpType(req *plugin.GenerateRequest, col *plugin.Column) string {
	switch req.Settings.Engine {
	case "sqlite":
		return sqliteType(col)
	case "mysql":
		return mysqlType(col)
	default:
		return "mixed"
	}
}

type goColumn struct {
	id      int
	docType string
	defVal  string
	*plugin.Column
}

func phpColumnsToStruct(req *plugin.GenerateRequest, name string, columns []goColumn, namer func(*plugin.Column, int) string) *ModelClass {
	gs := ModelClass{Name: name}
	idSeen := map[int]Field{}
	nameSeen := map[string]int{}
	for _, c := range columns {
		if _, ok := idSeen[c.id]; ok {
			continue
		}

		fieldName := memberName(namer(c.Column, c.id))
		if v := nameSeen[c.Name]; v > 0 {
			fieldName = fmt.Sprintf("%s_%d", fieldName, v+1)
		}

		field := Field{
			OriginalColumnName: c.Column.Name,
			ID:                 c.id,
			Name:               fieldName,
			Type:               makePhpTypeFromSqlcColumn(req, c.Column),
		}

		if c.docType != "" {
			field.Type = phpType{Name: c.docType}
		}

		if c.defVal != "" {
			field.Default = c.defVal
		}

		gs.Fields = append(gs.Fields, field)
		nameSeen[c.Name]++
		idSeen[c.id] = field
	}

	return &gs
}

func phpParamName(c *plugin.Column, number int) string {
	if c.Name != "" {
		return c.Name
	}

	return fmt.Sprintf("dollar_%d", number)
}

func phpColumnName(c *plugin.Column, pos int) string {
	if c.Name != "" {
		return c.Name
	}

	return fmt.Sprintf("column_%d", pos+1)
}

func BuildQueries(req *plugin.GenerateRequest, modelClasses []*ModelClass) ([]Query, []*ModelClass, error) {
	queries := make([]Query, 0, len(req.Queries))
	emitModelClasses := make([]*ModelClass, 0)

	for _, query := range req.Queries {
		if query.Name == "" || query.Cmd == "" {
			continue
		}

		if query.Cmd == metadata.CmdCopyFrom {
			return nil, nil, errors.New("support for CopyFrom in PHP is not implemented")
		}

		queryString := query.Text
		trimmedComments := make([]string, len(query.Comments))
		for i, c := range query.Comments {
			trimmedComments[i] = strings.TrimSpace(c)
		}

		queryStruct := Query{
			Cmd:          query.Cmd,
			ClassName:    strings.ToUpper(query.Name[:1]) + query.Name[1:],
			ConstantName: sdk.LowerTitle(query.Name),
			FieldName:    sdk.LowerTitle(query.Name) + "Stmt",
			MethodName:   sdk.LowerTitle(query.Name),
			SourceName:   query.Filename,
			SQL:          queryString,
			Comments:     trimmedComments,
		}

		var cols []goColumn
		overrides := parseSQLCParamComments(trimmedComments)
		for _, p := range query.Params {
			paramName := "$" + phpParamName(p.Column, int(p.Number))
			cols = append(cols, goColumn{
				id:      int(p.Number),
				Column:  p.Column,
				docType: overrides[paramName].typ,
				defVal:  overrides[paramName].def,
			})
		}

		params := phpColumnsToStruct(req, queryStruct.ClassName+"Bindings", cols, phpParamName)
		queryStruct.Arg = Params{ModelClass: params}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			queryStruct.Ret = QueryValue{
				Name: "results",
				Typ:  makePhpTypeFromSqlcColumn(req, c),
			}
		} else if len(query.Columns) > 1 {
			var gs *ModelClass

			for _, s := range modelClasses {
				if len(s.Fields) != len(query.Columns) {
					continue
				}

				same := true
				for i, f := range s.Fields {
					c := query.Columns[i]
					if f.Name != memberName(phpColumnName(c, i)) || f.Type != makePhpTypeFromSqlcColumn(req, c) || !sdk.SameTableName(c.Table, &s.Table, req.Catalog.DefaultSchema) {
						same = false
						break
					}
				}

				if same {
					gs = s
					break
				}
			}

			if gs == nil {
				var columns []goColumn
				for i, c := range query.Columns {
					columns = append(columns, goColumn{id: i, Column: c})
				}
				gs = phpColumnsToStruct(req, queryStruct.ClassName+"Row", columns, phpColumnName)
				emitModelClasses = append(emitModelClasses, gs)
			}

			queryStruct.Ret = QueryValue{
				Name:   "results",
				Struct: gs,
			}
		}

		queries = append(queries, queryStruct)
	}

	sort.Slice(queries, func(i, j int) bool { return queries[i].MethodName < queries[j].MethodName })
	return queries, emitModelClasses, nil
}
