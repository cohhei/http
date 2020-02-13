package http

import (
	"fmt"
	"io"
	"strings"
)

type Header map[string]string

func (h Header) writeTo(w io.Writer) (int, error) {
	l := len(h)
	if l == 0 {
		return fmt.Fprint(w, "\n")
	}

	header := make([]string, l)
	i := 0
	for k, v := range h {
		header[i] = fmt.Sprintf("%s: %s", k, v)
		i++
	}
	return fmt.Fprintf(w, "%s\n\n", strings.Join(header, "\n"))
}
