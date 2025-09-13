package core

import (
	"testing"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

func TestParams_isEmpty(t *testing.T) {
	mc := &ModelClass{Fields: []Field{}}
	p := Params{ModelClass: mc}
	if !p.isEmpty() {
		t.Errorf("Expected isEmpty to be true for empty fields")
	}

	mc.Fields = append(mc.Fields, Field{Name: "foo"})
	if p.isEmpty() {
		t.Errorf("Expected isEmpty to be false for non-empty fields")
	}
}

func TestParams_Args(t *testing.T) {
	mc := &ModelClass{Fields: []Field{{Name: "foo", Type: phpType{Name: "int"}}}}
	p := Params{ModelClass: mc}
	expected := "int $foo"
	if got := p.Args(); got != expected {
		t.Errorf("Args() = %q, want %q", got, expected)
	}
}

func TestParams_Bindings(t *testing.T) {
	mc := &ModelClass{Fields: []Field{{Name: "foo", Type: phpType{Name: "int"}}}}
	p := Params{ModelClass: mc}
	expected := "[$foo]"
	if got := p.Bindings(); got != expected {
		t.Errorf("Bindings() = %q, want %q", got, expected)
	}
}

func TestDataClassName(t *testing.T) {
	name := "foo_bar"
	expected := "FooBar"
	if got := dataClassName(name); got != expected {
		t.Errorf("dataClassName() = %q, want %q", got, expected)
	}
}

func TestMemberName(t *testing.T) {
	name := "foo_bar"
	expected := "fooBar"
	if got := memberName(name); got != expected {
		t.Errorf("memberName() = %q, want %q", got, expected)
	}
}

func TestPhpParamName(t *testing.T) {
	col := &plugin.Column{Name: "foo"}
	if got := phpParamName(col, 1); got != "foo" {
		t.Errorf("phpParamName() = %q", got)
	}
	col.Name = ""
	if got := phpParamName(col, 2); got != "dollar_2" {
		t.Errorf("phpParamName() = %q", got)
	}
}

func TestPhpColumnName(t *testing.T) {
	col := &plugin.Column{Name: "bar"}
	if got := phpColumnName(col, 0); got != "bar" {
		t.Errorf("phpColumnName() = %q", got)
	}

	col.Name = ""
	if got := phpColumnName(col, 1); got != "column_2" {
		t.Errorf("phpColumnName() = %q", got)
	}
}

func TestMakePhpTypeFromSqlcColumn(t *testing.T) {
	req := &plugin.GenerateRequest{Settings: &plugin.Settings{Engine: "sqlite"}}
	col := &plugin.Column{Type: &plugin.Identifier{Name: "INTEGER"}, NotNull: true}
	typ := makePhpTypeFromSqlcColumn(req, col)
	if typ.Name == "" {
		t.Errorf("Expected non-empty type name")
	}
}

func TestMapSqlColumnTypeToPhpType(t *testing.T) {
	req := &plugin.GenerateRequest{Settings: &plugin.Settings{Engine: "sqlite"}}
	col := &plugin.Column{Type: &plugin.Identifier{Name: "INTEGER"}}
	if got := mapSqlColumnTypeToPhpType(req, col); got == "" {
		t.Errorf("Expected non-empty type mapping")
	}
}

func TestPhpColumnsToStruct(t *testing.T) {
	req := &plugin.GenerateRequest{Settings: &plugin.Settings{Engine: "sqlite"}}
	columns := []goColumn{
		{
			id: 1,
			Column: &plugin.Column{
				Name:    "foo",
				Type:    &plugin.Identifier{Name: "INTEGER"},
				NotNull: true,
			},
		},
	}

	mc := phpColumnsToStruct(req, "TestStruct", columns, func(c *plugin.Column, i int) string { return c.Name })
	if mc.Name != "TestStruct" {
		t.Errorf("phpColumnsToStruct.Name = %q", mc.Name)
	}

	if len(mc.Fields) != 1 || mc.Fields[0].Name == "" {
		t.Errorf("phpColumnsToStruct.Fields = %+v", mc.Fields)
	}
}
