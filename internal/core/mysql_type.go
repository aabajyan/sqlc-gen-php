package core

import (
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func mysqlType(req *plugin.GenerateRequest, col *plugin.Column) (string, bool) {
	columnType := sdk.DataType(col.Type)

	switch columnType {

	case "varchar", "text", "char", "tinytext", "mediumtext", "longtext":
		return "string", false

	case "int", "integer", "smallint", "mediumint", "year", "bigint":
		return "int", false

	case "blob", "binary", "varbinary", "tinyblob", "mediumblob", "longblob":
		return "string", false

	case "double", "double precision", "real":
		return "float", false

	case "decimal", "dec", "fixed":
		return "string", false

	case "enum":
		return "string", false

	case "date", "datetime", "time":
		return "\\DateTimeImmutable", false

	case "timestamp":
		return "Instant", false

	case "boolean", "bool", "tinyint":
		return "boolean", false

	case "json":
		return "string", false

	case "any":
		return "mixed", false

	default:
		for _, schema := range req.Catalog.Schemas {
			for _, enum := range schema.Enums {
				if columnType == enum.Name {
					if schema.Name == req.Catalog.DefaultSchema {
						return dataClassName(enum.Name, req.Settings), true
					}
					return dataClassName(schema.Name+"_"+enum.Name, req.Settings), true
				}
			}
		}
		return "Any", false
	}
}
