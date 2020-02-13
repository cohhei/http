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
	data := "Hello World!"

	go func() {
		HandleFunc(path, func(w ResponseWriter, r *Request) {
			if r.Method != req.Method || r.Host != req.Host || r.Path != req.Path || r.Port != req.Port {
				t.Fatalf("\nwant:\t%+v\ngot:\t%+v\n", req, r)
			}
			w.Header()["Content-Type"] = "text/plain"
			w.WriteHeader(200)
			io.Copy(w, strings.NewReader(data))
		})
		if err := ListenAndServe(8881); err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(10 * time.Millisecond)

	c := NewClient()
	got, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	want := &Response{
		Status: 200,
		Body:   strings.NewReader(data),
	}
	if got.Status != want.Status {
		t.Fatalf("\nwant:\t%+v\ngot:\t%+v\n", want.Status, got.Status)
	}

	gb, _ := ioutil.ReadAll(got.Body)
	wb, _ := ioutil.ReadAll(want.Body)
	if string(gb) != string(wb) {
		t.Fatalf("\nwant:\t%s\ngot:\t%s\n", wb, gb)
	}
}
