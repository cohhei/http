package http

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

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