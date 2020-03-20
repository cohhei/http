package http

import (
	"fmt"
	"io"
	"strings"
)

// Client is an interface for the HTTP client.
type Client interface {
	Do(req *Request) (*Response, error)
}

type client struct {
	tcpClient TCPClient
	dnsClient DNSClient
}

// NewClient returns a new HTTP client.
func NewClient() Client {
	c := &client{
		tcpClient: NewTCPClient(),
		dnsClient: NewDNSClient(),
	}
	return c
}

// Do sends an HTTP request and returns an HTTP response as `io.ReadCloser`.
// The response should be closed.
func (c *client) Do(req *Request) (*Response, error) {
	if req.Port == 0 {
		req.Port = 80
	}

	ip, err := c.dnsClient.Resolve(req.Host)
	if err != nil {
		return nil, err
	}

	conn, err := c.tcpClient.Connect(ip, req.Port)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	write(conn, req)

	return parseResponse(conn)
}

func write(w io.Writer, req *Request) {
	var header []string
	for k, v := range req.Header {
		header = append(header, fmt.Sprintf("%s: %s", k, v))
	}
	if req.Path == "" {
		req.Path = "/"
	}
	const format = `%s %s HTTP/1.1
Host: %s:%d
%s

`
	fmt.Fprintf(w, format, req.Method, req.Path, req.Host, req.Port, strings.Join(header, "\n"))
	if req.Body != nil {
		io.Copy(w, req.Body)
	}
}
