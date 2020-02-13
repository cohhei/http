package http

import (
	"fmt"
	"io"
	"strings"
)

var DefaultServer = &Server{
	tcpClient: NewTCPClient(),
	router:    make(map[string]handler),
}

type Server struct {
	tcpClient TCPClient
	router
}

type handler func(w ResponseWriter, r *Request)

type router map[string]handler

func ListenAndServe(port int) error {
	return DefaultServer.ListenAndServe(port)
}

func (s *Server) ListenAndServe(port int) error {
	ip := [4]byte{127, 0, 0, 1}
	listener, err := s.tcpClient.Listen(ip, port)
	if err != nil {
		return err
	}
	defer listener.Close()

	ch := make(chan io.ReadWriteCloser)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			ch <- conn
		}
	}()

	for {
		conn := <-ch
		go func() {
			defer conn.Close()

			req, err := parseRequest(conn)
			if err != nil {
				io.Copy(conn, strings.NewReader("HTTP/1.1 404 Not Found\n\n404 Not Found"))
			}

			handler, exists := s.router[req.Path]
			if exists {
				w := &responseWriter{
					conn:   conn,
					header: make(Header),
				}
				handler(w, req)
			} else {
				io.Copy(conn, strings.NewReader("HTTP/1.1 404 Not Found\n\n404 Not Found"))
			}
		}()
	}
}

func HandleFunc(path string, handler handler) {
	DefaultServer.router[path] = handler
}

func (s *Server) HandleFunc(path string, handler handler) {
	s.router[path] = handler
}

// A ResponseWriter interface is used by an HTTP handler to
// construct an HTTP response.
type ResponseWriter interface {
	// Header returns the header map that will be sent by WriteHeader
	Header() Header

	// Write writes the data to the connection as part of an HTTP reply.
	Write([]byte) (int, error)

	// WriteHeader sends an HTTP response header with the provided
	// status code.
	//
	// If WriteHeader is not called explicitly, the first call to Write
	// will trigger an implicit WriteHeader(http.StatusOK).
	WriteHeader(statusCode int)
}

type responseWriter struct {
	conn        io.Writer // TCP connection
	header      Header    // Response header
	wroteHeader bool      // If it already wrote the HTTP header
}

func (w *responseWriter) Header() Header {
	return w.header
}

func (w *responseWriter) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(200)
	}

	return w.conn.Write(p)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	fmt.Fprintf(w.conn, "HTTP/1.1 %d %s\n", statusCode, statusText[statusCode])
	w.header.writeTo(w.conn)
	w.wroteHeader = true
}
