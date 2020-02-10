package http

import "io"

// Response is a struct for HTTP responses
type Response struct {
	Status  int               // e.g. 200
	Headers map[string]string // Headers is a map for HTTP headers.
	Body    io.ReadCloser     // Body is the HTTP response body
}
