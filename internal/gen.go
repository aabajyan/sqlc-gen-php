package kotlin

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"strings"
	_ "strings"
	"text/template"

	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
	"github.com/sqlc-dev/sqlc-gen-kotlin/internal/core"
)

//go:embed tmpl/models.tmpl
var modelsTemplate string

//go:embed tmpl/query_impl.tmpl
var queryImplTemplate string

//go:embed tmpl/query_interface.tmpl
var queryInterfaceTemplate string

func Offset(v int) int {
	return v + 1
}

func RemoveBlankLines(s string) string {
	skipNextSpace := false
	var lines []string
	for _, l := range strings.Split(s, "\n") {
		isSpace := len(strings.TrimSpace(l)) == 0
		if !isSpace || !skipNextSpace {
			lines = append(lines, l)
		}
		skipNextSpace = isSpace
	}
	o := strings.Join(lines, "\n")
	o += "\n"
	return o
}

func Generate(ctx context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	var conf core.Config
	if len(req.PluginOptions) > 0 {
		if err := json.Unmarshal(req.PluginOptions, &conf); err != nil {
			return nil, err
		}
	}

	structs := core.BuildDataClasses(req)
	queries, err := core.BuildQueries(req, structs)
	if err != nil {
		return nil, err
	}

	funcMap := template.FuncMap{
		"lowerTitle": sdk.LowerTitle,
		"comment":    sdk.DoubleSlashComment,
		"offset":     Offset,
	}

	modelsFile := template.Must(template.New("table").Funcs(funcMap).Parse(modelsTemplate))
	sqlFile := template.Must(template.New("table").Funcs(funcMap).Parse(queryImplTemplate))
	ifaceFile := template.Must(template.New("table").Funcs(funcMap).Parse(queryInterfaceTemplate))

	tctx := core.PhpTmplCtx{
		Settings:    req.Settings,
		Package:     conf.Package,
		Queries:     queries,
		DataClasses: structs,
		SqlcVersion: req.SqlcVersion,
	}

	output := map[string]string{}

	execute := func(name string, t *template.Template) error {
		tctx.SourceName = name
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		err := t.Execute(w, tctx)
		if err != nil {
			return err
		}
		err = w.Flush()
		if err != nil {
			return err
		}
		output[name] = RemoveBlankLines(b.String())
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
