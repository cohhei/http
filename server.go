package http

import (
	"bufio"
	"io"
	"strconv"
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

			req, err := parse(conn)
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

func parse(r io.Reader) (*Request, error) {
	// Read HTTP method and path
	buf := bufio.NewReader(r)
	line, _, err := buf.ReadLine()
	if err != nil {
		return nil, err
	}
	s := strings.Split(string(line), " ")
	req := &Request{
		Method:  s[0],
		Path:    s[1],
		Headers: make(map[string]string),
	}

	// Read headers
	for {
		line, _, err = buf.ReadLine()
		if err != nil {
			return nil, err
		}

		if len(line) == 0 {
			break
		}

		h := strings.Split(string(line), ": ")
		req.Headers[h[0]] = h[1]
		// Read host and port
		if h[0] == "Host" {
			uri := strings.Split(string(h[1]), ":")
			req.Host = uri[0]
			if len(uri) >= 2 {
				p, err := strconv.Atoi(uri[1])
				if err != nil {
					return nil, err
				}
				req.Port = p
			}
		}
	}

	req.Body = buf

	return req, nil
}

func HandleFunc(path string, handler handler) {
	DefaultServer.router[path] = handler
}

func (s *Server) HandleFunc(path string, handler handler) {
	s.router[path] = handler
}
