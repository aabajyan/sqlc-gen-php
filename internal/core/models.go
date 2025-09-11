package core

import "github.com/sqlc-dev/plugin-sdk-go/plugin"

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

type Field struct {
	ID                 int
	Name               string
	OriginalColumnName string
	Type               phpType
	Comment            string
	Default            string
	DocType            string
}

type ModelClass struct {
	Table   plugin.Identifier
	Name    string
	Fields  []Field
	Comment string
}

type QueryValue struct {
	Name   string
	Struct *ModelClass
	Typ    phpType
}

func (v QueryValue) IsStruct() bool {
	return v.Struct != nil
}

func (v QueryValue) Type() string {
	if v.Typ != (phpType{}) {
		return v.Typ.String()
	}

	if v.Struct != nil {
		return v.Struct.Name
	}

	panic("no type for QueryValue: " + v.Name)
}

func (v QueryValue) IsClass() bool {
	switch v.Typ.Name {
	case "int", "float", "string", "bool":
		return false
	default:
		return true
	}
}

type QueriesTmplCtx struct {
	Package     string
	Queries     []Query
	Settings    *plugin.Settings
	SqlcVersion string
	SourceName  string
}

type ModelsTmplCtx struct {
	Package     string
	ModelClass  *ModelClass
	SqlcVersion string
	SourceName  string
}

type phpType struct {
	Name     string
	IsArray  bool
	IsNull   bool
	DataType string
	Engine   string
}

func (t phpType) String() string {
	v := t.Name
	if t.IsArray {
		v = "array"
	} else if t.IsNull && t.Name != "mixed" {
		v = "?" + v
	}

	return v
}

func (t phpType) IsDateTimeImmutable() bool {
	return t.Name == "\\DateTimeImmutable"
}

func (t phpType) IsJSON() bool {
	return t.Name == "array"
}

func (t phpType) IsInt() bool {
	return t.Name == "int"
}

func (t phpType) IsFloat() bool {
	return t.Name == "float"
}

func (t phpType) IsString() bool {
	return t.Name == "string"
}
