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
