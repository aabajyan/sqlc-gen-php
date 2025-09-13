package core

import (
	"testing"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

func TestMysqlType(t *testing.T) {
	cases := []struct {
		name     string
		colType  string
		expected string
	}{
		{"varchar", "varchar", "string"},
		{"int", "int", "int"},
		{"blob", "blob", "string"},
		{"double", "double", "float"},
		{"decimal", "decimal", "string"},
		{"enum", "enum", "string"},
		{"date", "date", "string"},
		{"boolean", "boolean", "bool"},
		{"json", "json", "array"},
		{"any", "any", "mixed"},
		{"unknown", "unknown", "mixed"},
	}
	for _, tc := range cases {
		col := &plugin.Column{Type: &plugin.Identifier{Name: tc.colType}}
		if got := mysqlType(col); got != tc.expected {
			t.Errorf("mysqlType(%q) = %q, want %q", tc.colType, got, tc.expected)
		}
	}
}
