package core

import (
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func mysqlType(col *plugin.Column) string {
	columnType := sdk.DataType(col.Type)

	switch columnType {

	case "varchar", "text", "char", "tinytext", "mediumtext", "longtext":
		return "string"

	case "int", "integer", "smallint", "mediumint", "year", "bigint":
		return "int"

	case "blob", "binary", "varbinary", "tinyblob", "mediumblob", "longblob":
		return "string"

	case "double", "double precision", "real":
		return "float"

	case "decimal", "dec", "fixed":
		return "string"

	case "enum":
		return "string"

	case "date", "datetime", "time", "timestamp":
		return "\\DateTimeImmutable"

	case "boolean", "bool", "tinyint":
		return "boolean"

	case "json":
		return "string"

	case "any":
		return "mixed"

	default:
		return "mixed"
	}
}
