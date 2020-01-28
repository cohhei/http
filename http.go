package http

import (
	"fmt"
	"io"
	"net"
	"strings"
)

// Request is a struct for HTTP requests
type Request struct {
	// Method is a string for HTTP methods.
	Method string

	// Headers is a map for HTTP headers.
	Headers map[string]string

	// Host is a string for the hostname.
	Host string

	// Port is the port number.
	// If it is omitted, it will be 80.
	Port int

	// Path is the URL path.
	// If it is omitted, it will be '/'.
	Path string

	// Body is the HTTP request body.
	Body string
}

// Client is an interface for the HTTP client.
type Client interface {
	Do(req *Request) (io.ReadCloser, error)
}

type client struct {
	tcpClient TCPClient
}

func NewClient() Client {
	c := &client{
		tcpClient: NewTCPClient(),
	}
	return c
}

// Do sends an HTTP request and returns an HTTP response as `io.ReadCloser`.
// The response should be closed.
func (c *client) Do(req *Request) (io.ReadCloser, error) {
	if req.Port == 0 {
		req.Port = 80
	}

	addr, err := getIP(req.Host, req.Port)
	if err != nil {
		return nil, err
	}
	if err := c.tcpClient.Connect(addr); err != nil {
		return nil, err
	}

	data := rawRequest(req)
	if _, err := c.tcpClient.Send(data); err != nil {
		return nil, err
	}

	resp, err := c.tcpClient.GetReader()
	if err != nil {
		return nil, err
	}
	return &readCloser{resp, c.tcpClient}, nil
}

type readCloser struct {
	io.Reader
	TCPClient
}

func (r *readCloser) Close() error {
	return r.TCPClient.Close()
}

func getIP(hostname string, port int) (*Addr, error) {
	ipAddr, err := net.ResolveIPAddr("ip", hostname)
	if err != nil {
		return nil, err
	}

	var ip [4]byte
	copy(ip[:], ipAddr.IP.To4())
	if port == 0 {
		port = 80
	}
	addr := &Addr{
		IP:   ip,
		Port: port,
	}

	return addr, nil
}

func rawRequest(req *Request) []byte {
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

%s`
	return []byte(fmt.Sprintf(format, req.Method, req.Path, req.Host, req.Port, strings.Join(headers, "\n"), req.Body))
}
