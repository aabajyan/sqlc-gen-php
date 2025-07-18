package core

import (
	"testing"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

func TestSqliteType(t *testing.T) {
	cases := []struct {
		name     string
		colType  string
		expected string
	}{
		{"text", "text", "string"},
		{"integer", "integer", "int"},
		{"real", "real", "float"},
		{"blob", "blob", "string"},
		{"boolean", "boolean", "boolean"},
		{"date", "date", "\\DateTimeImmutable"},
		{"numeric", "numeric", "string"},
		{"json", "json", "string"},
		{"any", "any", "mixed"},
		{"unknown", "unknown", "mixed"},
	}
	for _, tc := range cases {
		col := &plugin.Column{Type: &plugin.Identifier{Name: tc.colType}}
		if got := sqliteType(col); got != tc.expected {
			t.Errorf("sqliteType(%q) = %q, want %q", tc.colType, got, tc.expected)
		}
	}
}
