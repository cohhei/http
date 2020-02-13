package http

import "io"

// Response is a struct for HTTP responses
type Response struct {
	Status int       // e.g. 200
	Header Header    // Header is a map for HTTP header.
	Body   io.Reader // Body is the HTTP response body
}
