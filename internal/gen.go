package generator

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"strings"
	_ "strings"
	"text/template"

	"github.com/lcarilla/sqlc-plugin-php-dbal/internal/core"
	"github.com/sqlc-dev/plugin-sdk-go/plugin"
	"github.com/sqlc-dev/plugin-sdk-go/sdk"
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

func executeTemplate(name string, t *template.Template, ctx interface{}, output map[string]string) error {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := t.Execute(w, ctx); err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return err
	}
	output[name] = RemoveBlankLines(b.String())
	return nil
}

func Generate(_ context.Context, req *plugin.GenerateRequest) (*plugin.GenerateResponse, error) {
	var conf core.Config
	if len(req.PluginOptions) > 0 {
		if err := json.Unmarshal(req.PluginOptions, &conf); err != nil {
			return nil, err
		}
	}

	modelClasses := core.BuildDataClasses(req)
	queries, emitModelClasses, err := core.BuildQueries(req, modelClasses)
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

	queryTemplateContext := core.QueriesTmplCtx{
		Settings:    req.Settings,
		Package:     conf.Package,
		Queries:     queries,
		SqlcVersion: req.SqlcVersion,
	}

	output := map[string]string{}

	if err := executeTemplate("Queries.php", ifaceFile, queryTemplateContext, output); err != nil {
		return nil, err
	}
	if err := executeTemplate("QueriesImpl.php", sqlFile, queryTemplateContext, output); err != nil {
		return nil, err
	}

	for i := range modelClasses {
		modelClass := &modelClasses[i]
		if err := executeTemplate(modelClass.Name+".php", modelsFile, &core.ModelsTmplCtx{
			Package:     conf.Package,
			SqlcVersion: req.SqlcVersion,
			ModelClass:  modelClass,
		}, output); err != nil {
			return nil, err
		}
	}

	for _, modelClass := range emitModelClasses {
		if err := executeTemplate(modelClass.Name+".php", modelsFile, &core.ModelsTmplCtx{
			Package:     conf.Package,
			SqlcVersion: req.SqlcVersion,
			ModelClass:  modelClass,
		}, output); err != nil {
			return nil, err
		}
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
