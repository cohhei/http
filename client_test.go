package http

import (
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
	defer resp.Close()

	buf := make([]byte, 15)
	if _, err := resp.Read(buf); err != nil {
		t.Fatal(err)
	}

	got := string(buf)
	want := "HTTP/1.1 200 OK"

	if got != want {
		t.Fatalf("expected:\t%s\nactual:\t%s\n", want, got)
	}
}

func TestRawRequest(t *testing.T) {
	req := &Request{
		Method: "GET",
		Headers: map[string]string{
			"Key":    "Value",
		},
		Host: "example.com",
		Port: 80,
		Body: "_",
	}
	got := string(rawRequest(req))
	want := `GET / HTTP/1.1
Host: example.com:80
Key: Value

_`

	if got != want {
		t.Fatalf("expected:\t%s\nactual:\t%s\n", want, got)
	}
}
