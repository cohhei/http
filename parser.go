package http

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

func parseRequest(r io.Reader) (*Request, error) {
	// Read HTTP method and path
	buf := bufio.NewReader(r)
	line, _, err := buf.ReadLine()
	if err != nil {
		return nil, err
	}
	s := strings.Split(string(line), " ")
	req := &Request{
		Method: s[0],
		Path:   s[1],
	}

	header, err := parseHeader(buf)
	if err != nil {
		return nil, err
	}
	req.Header = header

	// Read host and port
	if host, exists := header["Host"]; exists {
		uri := strings.Split(string(host), ":")
		req.Host = uri[0]
		req.Port = 80
		if len(uri) >= 2 {
			p, err := strconv.Atoi(uri[1])
			if err != nil {
				return nil, err
			}
			req.Port = p
		}
	}

	req.Body = buf

	return req, nil
}

func parseResponse(r io.Reader) (*Response, error) {
	// Read HTTP status code
	buf := bufio.NewReader(r)
	line, _, err := buf.ReadLine()
	if err != nil {
		return nil, err
	}
	s := strings.Split(string(line), " ")
	code, err := strconv.Atoi(s[1])
	if err != nil {
		return nil, err
	}
	resp := &Response{
		Status: code,
	}

	header, err := parseHeader(buf)
	if err != nil {
		return nil, err
	}
	resp.Header = header

	resp.Body = buf

	return resp, nil
}

func parseHeader(buf *bufio.Reader) (map[string]string, error) {
	header := make(map[string]string)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			return nil, err
		}
		if len(line) == 0 {
			break
		}
		h := strings.Split(string(line), ": ")
		header[h[0]] = h[1]
	}
	return header, nil
}
