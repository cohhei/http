package http

import (
	"io"
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
	Body io.Reader
}
