package kotlin

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	_ "strings"
	"text/template"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
	"github.com/sqlc-dev/sqlc-gen-kotlin/internal/core"
)

//go:embed tmpl/ktmodels.tmpl
var ktModelsTmpl string

//go:embed tmpl/ktsql.tmpl
var ktSqlTmpl string

//go:embed tmpl/ktiface.tmpl
var ktIfaceTmpl string

func Offset(v int) int {
	return v + 1
}

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	var conf core.Config
	if len(req.PluginOptions) > 0 {
		if err := json.Unmarshal(req.PluginOptions, &conf); err != nil {
			return nil, err
		}
	}

	structs := core.BuildDataClasses(conf, req)
	queries, err := core.BuildQueries(req, structs)
	if err != nil {
		return nil, err
	}

	i := &core.Importer{
		Settings:    req.Settings,
		DataClasses: structs,
		Queries:     queries,
	}

	funcMap := template.FuncMap{
		"lowerTitle": sdk.LowerTitle,
		"comment":    sdk.DoubleSlashComment,
		"imports":    i.Imports,
		"offset":     Offset,
	}

	modelsFile := template.Must(template.New("table").Funcs(funcMap).Parse(ktModelsTmpl))
	sqlFile := template.Must(template.New("table").Funcs(funcMap).Parse(ktSqlTmpl))
	ifaceFile := template.Must(template.New("table").Funcs(funcMap).Parse(ktIfaceTmpl))

	core.DefaultImporter = i

	tctx := core.KtTmplCtx{
		Settings:    req.Settings,
		Q:           `"""`,
		Package:     conf.Package,
		Queries:     queries,
		DataClasses: structs,
		SqlcVersion: req.SqlcVersion,
	}

	output := map[string]string{}

	execute := func(name string, t *template.Template) error {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		err := t.Execute(w, tctx)
		w.Flush()
		if err != nil {
			return err
		}
		output[name] = core.KtFormat(b.String())
		return nil
	}

	if err := execute("Models.php", modelsFile); err != nil {
		return nil, err
	}
	if err := execute("Queries.php", ifaceFile); err != nil {
		return nil, err
	}
	if err := execute("QueriesImpl.php", sqlFile); err != nil {
		return nil, err
	}

	resp := plugin.GenerateResponse{}

	for filename, code := range output {
		resp.Files = append(resp.Files, &plugin.File{
			Name:     filename,
			Contents: []byte(code),
		})
	}

	return &resp, nil
}
