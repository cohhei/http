package http

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// Client is an interface for the HTTP client.
type Client interface {
	Do(req *Request) (*Response, error)
}

type client struct {
	tcpClient TCPClient
}

// NewClient returns a new HTTP client.
func NewClient() Client {
	c := &client{
		tcpClient: NewTCPClient(),
	}
	return c
}

// Do sends an HTTP request and returns an HTTP response as `io.ReadCloser`.
// The response should be closed.
func (c *client) Do(req *Request) (*Response, error) {
	if req.Port == 0 {
		req.Port = 80
	}

	ip, port, err := getIP(req.Host, req.Port)
	if err != nil {
		return nil, err
	}
	conn, err := c.tcpClient.Connect(ip, port)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	write(conn, req)

	return parseResponse(conn)
}

func getIP(hostname string, port int) ([4]byte, int, error) {
	ipAddr, err := net.ResolveIPAddr("ip", hostname)
	if err != nil {
		return [4]byte{}, 0, err
	}

	var ip [4]byte
	copy(ip[:], ipAddr.IP.To4())
	if port == 0 {
		port = 80
	}

	return ip, port, nil
}

func write(w io.Writer, req *Request) {
	var headers []string
	for k, v := range req.Headers {
		headers = append(headers, fmt.Sprintf("%s: %s", k, v))
	}
	if req.Path == "" {
		req.Path = "/"
	}
	const format = `%s %s HTTP/1.1
Host: %s:%d
%s

`
	fmt.Fprintf(w, format, req.Method, req.Path, req.Host, req.Port, strings.Join(headers, "\n"))
	if req.Body != nil {
		io.Copy(w, req.Body)
	}
}
