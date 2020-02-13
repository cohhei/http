package http

import (
	"bytes"
	"testing"
)

func TestWriteTo(t *testing.T) {
	testCases := []struct {
		desc   string
		header Header
		want   string
	}{
		{
			desc:   "empty header",
			header: make(Header),
			want:   "\n",
		},
		{
			desc: "header",
			header: Header{
				"Content-Type":   "text/html",
				"Content-Length": "100",
			},
			want: "Content-Type: text/html\nContent-Length: 100\n\n",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})
			tC.header.writeTo(buf)

			got := buf.String()
			if got != tC.want {
				t.Fatalf("\nwant:\t%s\ngot:\t%s\n", tC.want, got)
			}
		})
	}
}
