package http

import (
	"bytes"
	"strings"
	"testing"
)

func TestDo(t *testing.T) {
	req := &Request{
		Method: "GET",
		Host:   "example.com",
		Path:   "/index.html",
	}

	c := NewClient()
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Status != 200 {
		t.Fatalf("\nwant:\t%d\ngot:\t%d\n", 200, resp.Status)
	}

	contentType := "text/html; charset=UTF-8"
	if resp.Header["Content-Type"] != contentType {
		t.Fatalf("\nwant:\t%s\ngot:\t%s\n", contentType, resp.Header["Content-Type"])
	}
}

func TestRawRequest(t *testing.T) {
	req := &Request{
		Method: "GET",
		Header: Header{
			"Key": "Value",
		},
		Host: "example.com",
		Port: 80,
		Body: strings.NewReader("_"),
	}

	buf := bytes.NewBuffer(nil)
	write(buf, req)
	got := buf.String()
	want := `GET / HTTP/1.1
Host: example.com:80
Key: Value

_`

	if got != want {
		t.Fatalf("expected:\t%s\nactual:\t%s\n", want, got)
	}
}
