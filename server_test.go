package http

import (
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	r := strings.NewReader(`GET /index.html HTTP/1.1
Host: localhost:8080
User-Agent: curl/7.64.1
Accept: */*

__`)

	got, err := parse(r)
	if err != nil {
		t.Fatal(err)
	}

	want := &Request{
		Method: "GET",
		Headers: map[string]string{
			"Host":       "localhost:8080",
			"User-Agent": "curl/7.64.1",
			"Accept":     "*/*",
		},
		Host: "localhost",
		Port: 8080,
		Path: "/index.html",
		Body: strings.NewReader("__"),
	}

	if got.Method != want.Method || got.Host != want.Host || got.Path != want.Path || got.Port != want.Port {
		t.Fatalf("\nwant:\t%+v\ngot:\t%+v\n", want, got)
	}

	if !reflect.DeepEqual(got.Headers, want.Headers) {
		t.Fatalf("\nwant:\t%+v\ngot:\t%+v\n", want.Headers, got.Headers)
	}

	gb, _ := ioutil.ReadAll(got.Body)
	wb, _ := ioutil.ReadAll(want.Body)
	if string(gb) != string(wb) {
		t.Fatalf("\nwant:\t%s\ngot:\t%s\n", wb, gb)
	}
}

func TestListenAndServe(t *testing.T) {
	path := "/path/to/content"
	req := &Request{
		Method: "GET",
		Host:   "127.0.0.1",
		Port:   8888,
		Path:   path,
	}
	data := "HTTP/1.1 200 OK\n\nHello World!"

	go func() {
		HandleFunc(path, func(w io.Writer, r *Request) {
			if r.Method != req.Method || r.Host != req.Host || r.Path != req.Path || r.Port != req.Port {
				t.Fatalf("\nwant:\t%+v\ngot:\t%+v\n", req, r)
			}
			io.Copy(w, strings.NewReader(data))
			// w.Write([]byte(data))
		})
		if err := ListenAndServe(8888); err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(10 * time.Millisecond)

	c := NewClient()
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadAll(resp)
	if err != nil {
		t.Fatal(err)
	}

	got := strings.TrimSpace(string(b))
	want := data
	if got != want {
		t.Fatalf("\nwant:\t%+v\ngot:\t%+v\n", []byte(want), []byte(got))
	}
}
