package http

import (
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

type handler func(w io.Writer, r *Request)

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
				handler(conn, req)
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
