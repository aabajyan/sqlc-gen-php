package core

import (
	"strings"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func sqliteType(col *plugin.Column) string {
	columnType := strings.ToLower(sdk.DataType(col.Type))

	switch columnType {
	case "text", "varchar", "char", "clob":
		return "string"
	case "integer", "int", "bigint", "smallint", "tinyint":
		return "int"
	case "real", "double", "float":
		return "float"
	case "blob":
		return "string"
	case "boolean":
		return "boolean"
	case "date", "datetime":
		return "\\DateTimeImmutable"
	case "numeric":
		return "string"
	case "json":
		return "array"
	case "any":
		return "mixed"
	default:
		return "mixed"
	}
}
