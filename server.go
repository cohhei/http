package http

import "io"

var DefaultServer = &Server{
	tcpClient: NewTCPClient(),
	router:    make(map[string]handler),
}

type Server struct {
	tcpClient TCPClient
	router
}

type handler func(w io.Writer, r io.Reader)

type router map[string]handler

func ListenAndServe(port int) error {
	return DefaultServer.ListenAndServe(port)
}

func (s *Server) ListenAndServe(port int) error {
	ip := [4]byte{127, 0, 0, 1}
	conn, err := s.tcpClient.Listen(ip, port)
	if err != nil {
		return err
	}
	defer conn.Close()

	s.router["/"](conn, conn)

	if err := conn.Close(); err != nil {
		return err
	}
	return nil
}

func HandleFunc(path string, handler handler) {
	DefaultServer.router[path] = handler
}

func (s *Server) HandleFunc(path string, handler handler) {
	s.router[path] = handler
}
