package core

import (
	"log"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
)

func postgresType(req *plugin.GenerateRequest, col *plugin.Column) (string, bool) {
	columnType := sdk.DataType(col.Type)

	switch columnType {
	case "serial", "pg_catalog.serial4":
		return "int", false

	case "bigserial", "pg_catalog.serial8":
		return "int", false

	case "smallserial", "pg_catalog.serial2":
		return "int", false

	case "integer", "int", "int4", "pg_catalog.int4":
		return "int", false

	case "bigint", "pg_catalog.int8":
		return "int", false

	case "smallint", "pg_catalog.int2":
		return "int", false

	case "float", "double precision", "pg_catalog.float8":
		return "float", false

	case "real", "pg_catalog.float4":
		return "float", false

	case "pg_catalog.numeric":
		return "float", false

	case "bool", "pg_catalog.bool":
		return "boolean", false

	case "jsonb":
		// TODO: support json and byte types
		return "string", false

	case "bytea", "blob", "pg_catalog.bytea":
		return "string", false

	case "date":
		// Date and time mappings from https://jdbc.postgresql.org/documentation/head/java8-date-time.html
		return "\\DateTimeImmutable", false

	case "pg_catalog.time", "pg_catalog.timetz":
		return "\\DateTimeImmutable", false

	case "pg_catalog.timestamp":
		return "\\DateTimeImmutable", false

	case "pg_catalog.timestamptz", "timestamptz":
		// TODO
		return "\\DateTimeImmutable", false

	case "text", "pg_catalog.varchar", "pg_catalog.bpchar", "string":
		return "string", false

	case "uuid":
		return "Uuid", false

	case "inet":
		// TODO
		return "string", false

	case "void":
		// TODO
		// A void value always returns NULL. Since there is no built-in NULL
		// value into the SQL package, we'll use sql.NullBool
		return "sql.NullBool", false

	case "any":
		// TODO
		return "Any", false

	default:
		for _, schema := range req.Catalog.Schemas {
			if schema.Name == "pg_catalog" || schema.Name == "information_schema" {
				continue
			}
			for _, enum := range schema.Enums {
				if columnType == enum.Name {
					if schema.Name == req.Catalog.DefaultSchema {
						return dataClassName(enum.Name), true
					}
					return dataClassName(schema.Name + "_" + enum.Name), true
				}
			}
		}
		log.Printf("unknown PostgreSQL type: %s\n", columnType)
		return "Any", false
	}
}
