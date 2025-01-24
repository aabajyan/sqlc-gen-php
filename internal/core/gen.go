package core

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/sqlc-dev/plugin-sdk-go/metadata"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

type Constant struct {
	Name  string
	Type  string
	Value string
}

type Enum struct {
	Name      string
	Comment   string
	Constants []Constant
}

type Field struct {
	ID      int
	Name    string
	Type    ktType
	Comment string
}

type Struct struct {
	Table   plugin.Identifier
	Name    string
	Fields  []Field
	Comment string
}

type QueryValue struct {
	Emit   bool
	Name   string
	Struct *Struct
	Typ    ktType
}

func (v QueryValue) EmitStruct() bool {
	return v.Emit
}

func (v QueryValue) IsStruct() bool {
	return v.Struct != nil
}

func (v QueryValue) isEmpty() bool {
	return v.Typ == (ktType{}) && v.Name == "" && v.Struct == nil
}

func (v QueryValue) Type() string {
	if v.Typ != (ktType{}) {
		return v.Typ.String()
	}
	if v.Struct != nil {
		return v.Struct.Name
	}
	panic("no type for QueryValue: " + v.Name)
}

func jdbcSet(t ktType, idx int, name string) string {
	if t.IsEnum && t.IsArray {
		return fmt.Sprintf(`stmt.setArray(%d, conn.createArrayOf("%s", %s.map { v -> v.value }.toTypedArray()))`, idx, t.DataType, name)
	}
	if t.IsArray {
		return fmt.Sprintf(`stmt.setArray(%d, conn.createArrayOf("%s", %s.toTypedArray()))`, idx, t.DataType, name)
	}
	if t.IsTime() {
		return fmt.Sprintf("stmt.setObject(%d, %s)", idx, name)
	}
	if t.IsInstant() {
		return fmt.Sprintf("stmt.setTimestamp(%d, Timestamp.from(%s))", idx, name)
	}
	if t.IsUUID() {
		return fmt.Sprintf("stmt.setObject(%d, %s)", idx, name)
	}
	return fmt.Sprintf("'%d' => $%s,", idx, name)
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
	for i, f := range v.Struct.Fields {
		out = append(out, jdbcSet(f.Type, i+1, f.Name))
	}
	return indent(strings.Join(out, "\n"), 10, 0)
}

func jdbcGet(t ktType, idx int, name string) string {
	if t.IsEnum && t.IsArray {
		return fmt.Sprintf(`(results.getArray(%d).array as Array<String>).map { v -> %s.lookup(v)!! }.toList()`, idx, t.Name)
	}
	if t.IsEnum {
		return fmt.Sprintf("%s.lookup(results.getString(%d))!!", t.Name, idx)
	}
	if t.IsArray {
		return fmt.Sprintf(`(results.getArray(%d).array as Array<%s>).toList()`, idx, t.Name)
	}
	if t.IsTime() {
		return fmt.Sprintf(`results.getObject(%d, %s::class.java)`, idx, t.Name)
	}
	if t.IsInstant() {
		return fmt.Sprintf(`results.getTimestamp(%d).toInstant()`, idx)
	}
	if t.IsUUID() {
		return fmt.Sprintf(`$row["%s"] == null ? null : Uuid::fromString($row["%s"])`, name, name)
	}
	if t.IsBigDecimal() {
		return fmt.Sprintf(`results.getBigDecimal(%d)`, idx)
	}
	return fmt.Sprintf(`$row["%s"]`, name)
}

func (v QueryValue) ResultSet() string {
	var out []string
	for i, f := range v.Struct.Fields {
		out = append(out, jdbcGet(f.Type, i+1, f.Name))
	}
	ret := indent(strings.Join(out, ",\n"), 4, -1)
	return ret
}

func indent(s string, n int, firstIndent int) string {
	lines := strings.Split(s, "\n")
	buf := bytes.NewBuffer(nil)
	for i, l := range lines {
		indent := n
		if i == 0 && firstIndent != -1 {
			indent = firstIndent
		}
		if i != 0 {
			buf.WriteRune('\n')
		}
		for i := 0; i < indent; i++ {
			buf.WriteRune(' ')
		}
		buf.WriteString(l)
	}
	return buf.String()
}

// A struct used to generate methods and fields on the Queries struct
type Query struct {
	ClassName    string
	Cmd          string
	Comments     []string
	MethodName   string
	FieldName    string
	ConstantName string
	SQL          string
	SourceName   string
	Ret          QueryValue
	Arg          Params
}

func dataClassName(name string, settings *plugin.Settings) string {
	out := ""
	for _, p := range strings.Split(name, "_") {
		out += strings.Title(p)
	}
	return out
}

func memberName(name string, settings *plugin.Settings) string {
	return sdk.LowerTitle(dataClassName(name, settings))
}

func BuildDataClasses(conf Config, req *plugin.GenerateRequest) []Struct {
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
			structName := dataClassName(tableName, req.Settings)
			s := Struct{
				Table:   plugin.Identifier{Schema: schema.Name, Name: table.Rel.Name},
				Name:    structName,
				Comment: table.Comment,
			}
			for _, column := range table.Columns {
				s.Fields = append(s.Fields, Field{
					Name:    memberName(column.Name, req.Settings),
					Type:    makeType(req, column),
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

type ktType struct {
	Name     string
	IsEnum   bool
	IsArray  bool
	IsNull   bool
	DataType string
	Engine   string
}

func (t ktType) String() string {
	v := t.Name
	if t.IsArray {
		v = fmt.Sprintf("List<%s>", v)
	} else if t.IsNull {
		v = "?" + v
	}
	return v
}

func (t ktType) IsTime() bool {
	return t.Name == "LocalDate" || t.Name == "LocalDateTime" || t.Name == "LocalTime" || t.Name == "OffsetDateTime"
}

func (t ktType) IsInstant() bool {
	return t.Name == "Instant"
}

func (t ktType) IsUUID() bool {
	return t.Name == "UUID"
}

func (t ktType) IsBigDecimal() bool {
	return t.Name == "java.math.BigDecimal"
}

func makeType(req *plugin.GenerateRequest, col *plugin.Column) ktType {
	typ, isEnum := ktInnerType(req, col)
	return ktType{
		Name:     typ,
		IsEnum:   isEnum,
		IsArray:  col.IsArray,
		IsNull:   !col.NotNull,
		DataType: sdk.DataType(col.Type),
		Engine:   req.Settings.Engine,
	}
}

func ktInnerType(req *plugin.GenerateRequest, col *plugin.Column) (string, bool) {
	// TODO: Extend the engine interface to handle types
	switch req.Settings.Engine {
	case "mysql":
		return mysqlType(req, col)
	case "postgresql":
		return postgresType(req, col)
	default:
		return "Any", false
	}
}

type goColumn struct {
	id int
	*plugin.Column
}

func ktColumnsToStruct(req *plugin.GenerateRequest, name string, columns []goColumn, namer func(*plugin.Column, int) string) *Struct {
	gs := Struct{
		Name: name,
	}
	idSeen := map[int]Field{}
	nameSeen := map[string]int{}
	for _, c := range columns {
		if _, ok := idSeen[c.id]; ok {
			continue
		}
		fieldName := memberName(namer(c.Column, c.id), req.Settings)
		if v := nameSeen[c.Name]; v > 0 {
			fieldName = fmt.Sprintf("%s_%d", fieldName, v+1)
		}
		field := Field{
			ID:   c.id,
			Name: fieldName,
			Type: makeType(req, c.Column),
		}
		gs.Fields = append(gs.Fields, field)
		nameSeen[c.Name]++
		idSeen[c.id] = field
	}
	return &gs
}

func ktArgName(name string) string {
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

func ktParamName(c *plugin.Column, number int) string {
	if c.Name != "" {
		return ktArgName(c.Name)
	}
	return fmt.Sprintf("dollar_%d", number)
}

func ktColumnName(c *plugin.Column, pos int) string {
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
			ClassName:    strings.Title(query.Name),
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
		params := ktColumnsToStruct(req, gq.ClassName+"Bindings", cols, ktParamName)
		gq.Arg = Params{
			Struct: params,
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = QueryValue{
				Name: "results",
				Typ:  makeType(req, c),
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
					sameName := f.Name == memberName(ktColumnName(c, i), req.Settings)
					sameType := f.Type == makeType(req, c)
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
				gs = ktColumnsToStruct(req, gq.ClassName+"Row", columns, ktColumnName)
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

type KtTmplCtx struct {
	Package     string
	DataClasses []Struct
	Queries     []Query
	Settings    *plugin.Settings
	SqlcVersion string

	// TODO: Race conditions
	SourceName string

	EmitJSONTags        bool
	EmitPreparedQueries bool
	EmitInterface       bool
}
