package core

import (
	"sort"
	"strings"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
)

type Importer struct {
	Settings    *plugin.Settings
	DataClasses []Struct
	Queries     []Query
}

func (i *Importer) usesType(typ string) bool {
	for _, strct := range i.DataClasses {
		for _, f := range strct.Fields {
			if f.Type.Name == typ {
				return true
			}
		}
	}
	return false
}

func (i *Importer) Imports(filename string) [][]string {
	switch filename {
	case "Models.kt":
		return i.modelImports()
	case "Querier.kt":
		return i.interfaceImports()
	default:
		return i.queryImports(filename)
	}
}

func (i *Importer) interfaceImports() [][]string {
	uses := func(name string) bool {
		for _, q := range i.Queries {
			if !q.Ret.isEmpty() {
				if strings.HasPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				for _, f := range q.Arg.Struct.Fields {
					if strings.HasPrefix(f.Type.Name, name) {
						return true
					}
				}
			}
		}
		return false
	}

	std := stdImports(uses)
	stds := make([]string, 0, len(std))
	for s := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds}
}

func (i *Importer) modelImports() [][]string {
	std := make(map[string]struct{})
	if i.usesType("UUID") {
		std["Symfony\\Component\\Uid\\Uuid"] = struct{}{}
	}

	stds := make([]string, 0, len(std))
	for s := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds}
}

func stdImports(uses func(name string) bool) map[string]struct{} {
	std := map[string]struct{}{}
	if uses("UUID") {
		std["Symfony\\Component\\Uid\\Uuid"] = struct{}{}
	}

	return std
}

func (i *Importer) queryImports(filename string) [][]string {
	uses := func(name string) bool {
		for _, q := range i.Queries {
			if !q.Ret.isEmpty() {
				if q.Ret.Struct != nil {
					for _, f := range q.Ret.Struct.Fields {
						if f.Type.Name == name {
							return true
						}
					}
				}
				if q.Ret.Type() == name {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				for _, f := range q.Arg.Struct.Fields {
					if f.Type.Name == name {
						return true
					}
				}
			}
		}
		return false
	}

	hasEnum := func() bool {
		for _, q := range i.Queries {
			if !q.Arg.isEmpty() {
				for _, f := range q.Arg.Struct.Fields {
					if f.Type.IsEnum {
						return true
					}
				}
			}
		}
		return false
	}

	std := stdImports(uses)
	if hasEnum() && i.Settings.Engine == "postgresql" {
		std["java.sql.Types"] = struct{}{}
	}

	stds := make([]string, 0, len(std))
	for s := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	return [][]string{stds}
}
