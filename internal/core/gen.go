package core

import (
	"errors"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"sort"
	"strings"

	"github.com/sqlc-dev/plugin-sdk-go/metadata"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func dbalParameter(name string) string {
	return fmt.Sprintf("$%s,", name)
}

type Params struct {
	Struct *Struct
}

func (v Params) isEmpty() bool {
	return len(v.Struct.Fields) == 0
}

func (v Params) Args() string {
	if v.isEmpty() {
		return ""
	}
	var out []string
	fields := v.Struct.Fields
	for _, f := range fields {
		out = append(out, f.Type.String()+" $"+f.Name)
	}
	if len(out) < 3 {
		return strings.Join(out, ", ")
	}
	return "\n" + indent(strings.Join(out, ",\n"), 6, -1)
}

func (v Params) Bindings() string {
	if v.isEmpty() {
		return ""
	}
	var out []string
	for _, f := range v.Struct.Fields {
		out = append(out, dbalParameter(f.Name))
	}
	return indent(strings.Join(out, "\n"), 10, 0)
}

func dbalType(t phpType) string {
	if t.IsInt() {
		if t.IsArray {
			return "ArrayParameterType::INTEGER"
		}
		return "ParameterType::INTEGER"
	}
	if t.IsDateTimeImmutable() {
		return "Types::DATE_IMMUTABLE"
	}
	if t.IsString() {
		if t.IsArray {
			return "ArrayParameterType::STRING"
		}
		return "ParameterType::STRING"
	}
	return "ParameterType::STRING"
}

func (v Params) DBALTypes() string {
	if v.isEmpty() {
		return ""
	}
	var out []string
	for _, f := range v.Struct.Fields {
		out = append(out, dbalType(f.Type)+",")
	}
	return indent(strings.Join(out, "\n"), 10, 0)
}

func dbalRowMapping(t phpType, name string) string {
	if t.IsDateTimeImmutable() {
		return fmt.Sprintf(`$row["%s"] == null ? null : new \DateTimeImmutable($row["%s"])`, name, name)
	}
	return fmt.Sprintf(`$row["%s"]`, name)
}

func (v QueryValue) ResultSet() string {
	var out []string
	for _, f := range v.Struct.Fields {
		out = append(out, dbalRowMapping(f.Type, f.Name))
	}
	ret := indent(strings.Join(out, ",\n"), 4, -1)
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

func BuildDataClasses(req *plugin.GenerateRequest) []Struct {
	var structs []Struct
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
			s := Struct{
				Table:   plugin.Identifier{Schema: schema.Name, Name: table.Rel.Name},
				Name:    structName,
				Comment: table.Comment,
			}
			for _, column := range table.Columns {
				s.Fields = append(s.Fields, Field{
					Name:    memberName(column.Name),
					Type:    makePhpTypeFromSqlcColumn(req, column),
					Comment: column.Comment,
				})
			}
			structs = append(structs, s)
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
	case "mysql":
		return mysqlType(col)
	default:
		return "Any"
	}
}

type goColumn struct {
	id int
	*plugin.Column
}

func phpColumnsToStruct(req *plugin.GenerateRequest, name string, columns []goColumn, namer func(*plugin.Column, int) string) *Struct {
	gs := Struct{
		Name: name,
	}
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
			ID:   c.id,
			Name: fieldName,
			Type: makePhpTypeFromSqlcColumn(req, c.Column),
		}
		gs.Fields = append(gs.Fields, field)
		nameSeen[c.Name]++
		idSeen[c.id] = field
	}
	return &gs
}

func phpFunctionArgumentName(name string) string {
	out := ""
	for i, p := range strings.Split(name, "_") {
		if i == 0 {
			out += strings.ToLower(p)
		} else {
			out += strings.Title(p)
		}
	}
	return out
}

func phpParamName(c *plugin.Column, number int) string {
	if c.Name != "" {
		return phpFunctionArgumentName(c.Name)
	}
	return fmt.Sprintf("dollar_%d", number)
}

func phpColumnName(c *plugin.Column, pos int) string {
	if c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("column_%d", pos+1)
}

func BuildQueries(req *plugin.GenerateRequest, structs []Struct) ([]Query, error) {
	qs := make([]Query, 0, len(req.Queries))
	for _, query := range req.Queries {
		if query.Name == "" {
			continue
		}
		if query.Cmd == "" {
			continue
		}
		if query.Cmd == metadata.CmdCopyFrom {
			return nil, errors.New("support for CopyFrom in PHP is not implemented")
		}

		ql := query.Text
		gq := Query{
			Cmd:          query.Cmd,
			ClassName:    cases.Title(language.English).String(query.Name),
			ConstantName: sdk.LowerTitle(query.Name),
			FieldName:    sdk.LowerTitle(query.Name) + "Stmt",
			MethodName:   sdk.LowerTitle(query.Name),
			SourceName:   query.Filename,
			SQL:          ql,
			Comments:     query.Comments,
		}

		var cols []goColumn
		for _, p := range query.Params {
			cols = append(cols, goColumn{
				id:     int(p.Number),
				Column: p.Column,
			})
		}
		params := phpColumnsToStruct(req, gq.ClassName+"Bindings", cols, phpParamName)
		gq.Arg = Params{
			Struct: params,
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = QueryValue{
				Name: "results",
				Typ:  makePhpTypeFromSqlcColumn(req, c),
			}
		} else if len(query.Columns) > 1 {
			var gs *Struct
			var emit bool

			for _, s := range structs {
				if len(s.Fields) != len(query.Columns) {
					continue
				}
				same := true
				for i, f := range s.Fields {
					c := query.Columns[i]
					sameName := f.Name == memberName(phpColumnName(c, i))
					sameType := f.Type == makePhpTypeFromSqlcColumn(req, c)
					sameTable := sdk.SameTableName(c.Table, &s.Table, req.Catalog.DefaultSchema)

					if !sameName || !sameType || !sameTable {
						same = false
					}
				}
				if same {
					gs = &s
					break
				}
			}

			if gs == nil {
				var columns []goColumn
				for i, c := range query.Columns {
					columns = append(columns, goColumn{
						id:     i,
						Column: c,
					})
				}
				gs = phpColumnsToStruct(req, gq.ClassName+"Row", columns, phpColumnName)
				emit = true
			}
			gq.Ret = QueryValue{
				Emit:   emit,
				Name:   "results",
				Struct: gs,
			}
		}

		qs = append(qs, gq)
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].MethodName < qs[j].MethodName })
	return qs, nil
}
