package http

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

var requests = []struct {
	name string
	raw  string
	want *Request
}{
	{
		name: "GET",
		raw:  fmt.Sprintf("GET / HTTP/1.1\n%s", headers[0].raw),
		want: &Request{
			Method:  "GET",
			Headers: headers[0].want,
			Host:    "example.com",
			Port:    8080,
			Path:    "/",
			Body:    strings.NewReader(""),
		},
	},
	{
		name: "POST",
		raw:  fmt.Sprintf("POST /new HTTP/1.1\n%sfield1=value1&field2=value2", headers[1].raw),
		want: &Request{
			Method:  "POST",
			Headers: headers[1].want,
			Host:    "example.com",
			Port:    80,
			Path:    "/new",
			Body:    strings.NewReader("field1=value1&field2=value2"),
		},
	},
}

var responses = []struct {
	name string
	raw  string
	want *Response
}{
	{
		name: "Status OK",
		raw:  fmt.Sprintf("HTTP/1.1 200 OK\n%s", headers[2].raw),
		want: &Response{
			Status:  200,
			Headers: headers[2].want,
			Body:    strings.NewReader(""),
		},
	},
	{
		name: "Not Found",
		raw:  fmt.Sprintf("HTTP/1.1 404 Not Found\n%s404 Not Found", headers[3].raw),
		want: &Response{
			Status:  404,
			Headers: headers[3].want,
			Body:    strings.NewReader("404 Not Found"),
		},
	},
}

var headers = []struct {
	name string
	raw  string
	want map[string]string
}{
	{
		"GET",
		"Host: example.com:8080\nUser-Agent: cohhei/http\nAccept: */*\n\n",
		map[string]string{
			"Host":       "example.com:8080",
			"User-Agent": "cohhei/http",
			"Accept":     "*/*",
		},
	},
	{
		"POST",
		"Host: example.com\nContent-Type: application/x-www-form-urlencoded\nContent-Length: 27\nUser-Agent: cohhei/http\n\n",
		map[string]string{
			"Host":           "example.com",
			"Content-Type":   "application/x-www-form-urlencoded",
			"Content-Length": "27",
			"User-Agent":     "cohhei/http",
		},
	},
	{
		"JSON",
		"Content-Type: application/json\nContent-Length: 2\n\n",
		map[string]string{
			"Content-Type":   "application/json",
			"Content-Length": "2",
		},
	},
	{
		"404",
		"Content-Type: text/plain\nContent-Length: 13\n\n",
		map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": "13",
		},
	},
}

func TestParseHeader(t *testing.T) {
	for _, header := range headers {
		t.Run(header.name, func(t *testing.T) {
			got, err := parseHeaders(bufio.NewReader(strings.NewReader(header.raw)))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, header.want) {
				t.Fatalf("\nwant:\t%v\ngot:\t%v\n", header.want, got)
			}
		})
	}
}

func TestParseRequest(t *testing.T) {
	for _, req := range requests {
		t.Run(req.name, func(t *testing.T) {
			got, err := parseRequest(strings.NewReader(req.raw))
			if err != nil {
				t.Fatal(err)
			}

			if got.Method != req.want.Method || got.Host != req.want.Host || got.Path != req.want.Path || got.Port != req.want.Port {
				t.Fatalf("\nwant:\t%+v\ngot:\t%+v\n", req.want, got)
			}

			if !reflect.DeepEqual(got.Headers, req.want.Headers) {
				t.Fatalf("\nwant:\t%+v\ngot:\t%+v\n", req.want.Headers, got.Headers)
			}

			gb, _ := ioutil.ReadAll(got.Body)
			wb, _ := ioutil.ReadAll(req.want.Body)
			if string(gb) != string(wb) {
				t.Fatalf("\nwant:\t%s\ngot:\t%s\n", wb, gb)
			}
		})
	}
}

func TestParseResponse(t *testing.T) {
	for _, resp := range responses {
		t.Run(resp.name, func(t *testing.T) {
			got, err := parseResponse(strings.NewReader(resp.raw))
			if err != nil {
				t.Fatal(err)
			}

			if got.Status != resp.want.Status {
				t.Fatalf("\nwant:\t%+v\ngot:\t%+v\n", resp.want.Status, got.Status)
			}

			if !reflect.DeepEqual(got.Headers, resp.want.Headers) {
				t.Fatalf("\nwant:\t%+v\ngot:\t%+v\n", resp.want.Headers, got.Headers)
			}

			gb, _ := ioutil.ReadAll(got.Body)
			wb, _ := ioutil.ReadAll(resp.want.Body)
			if string(gb) != string(wb) {
				t.Fatalf("\nwant:\t%s\ngot:\t%s\n", wb, gb)
			}
		})
	}
}
