package core

import (
	"bytes"
	"strings"
)

func indent(s string, n int, firstIndent int) string {
	lines := strings.Split(s, "\n")
	buf := bytes.NewBuffer(nil)
	for i, l := range lines {
		indent := n
		if i == 0 && firstIndent != -1 {
			indent = firstIndent
		}
		if i != 0 {
			buf.WriteRune('\n')
		}
		for i := 0; i < indent; i++ {
			buf.WriteRune(' ')
		}
		buf.WriteString(l)
	}
	return buf.String()
}

type paramOverride struct{ typ, def string }

// FIXME: This is very barebones and only supports simple format
func parseSQLCParamComments(comments []string) map[string]paramOverride {
	out := map[string]paramOverride{}
	for _, c := range comments {
		line := strings.TrimSpace(c)
		if !strings.HasPrefix(line, "@sqlc-param") {
			continue
		}

		toks := strings.Fields(line)
		if len(toks) < 3 {
			continue
		}

		typ := toks[1]
		nameToken := toks[2]
		name := nameToken
		def := ""
		if eq := strings.Index(nameToken, "="); eq != -1 {
			name = nameToken[:eq]
			def = strings.TrimSpace(nameToken[eq+1:])
		}

		out[name] = paramOverride{typ: typ, def: def}
	}
	return out
}
