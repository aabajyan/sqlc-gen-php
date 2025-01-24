package core

type Config struct {
	Package                     string   `json:"package"`
	InflectionExcludeTableNames []string `json:"inflection_exclude_table_names"`
}
