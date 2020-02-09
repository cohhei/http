package http

import (
	"bytes"
	"io"
	"syscall"
)

// TCPClient is an interface for the TPC client
type TCPClient interface {
	// Connect creates a TCP connection to the given IP address.
	// The returned connection should be closed.
	Connect(ip [4]byte, port int) (connection io.ReadWriteCloser, err error)

	// Listen listens for TCP connections.
	// This method creates a new socket and binds it to the given address by using syscall.Bind().
	// In addition, calls syscall.Listen() and syscall.Accept().
	// The process will be blocked until a request comes.
	// The returned connection should be closed.
	Listen(ip [4]byte, port int) (connection io.ReadWriteCloser, err error)
}

// NewTCPClient creates a client
func NewTCPClient() TCPClient {
	return &tcpClient{}
}

type tcpClient struct{}

func (c *tcpClient) Connect(ip [4]byte, port int) (io.ReadWriteCloser, error) {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}

	sa := sockaddr(ip, port)
	if err := syscall.Connect(socket, sa); err != nil {
		return nil, err
	}

	return &tcpConnection{socket, sa}, nil
}

func (c *tcpClient) Listen(ip [4]byte, port int) (io.ReadWriteCloser, error) {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.Close(socket)

	if err := syscall.Bind(socket, sockaddr(ip, port)); err != nil {
		return nil, err
	}

	if err := syscall.Listen(socket, 5); err != nil {
		return nil, err
	}

	nfd, sa, err := syscall.Accept(socket)
	if err != nil {
		return nil, err
	}

	return &tcpConnection{nfd, sa}, nil
}

type tcpConnection struct {
	socket int
	sa     syscall.Sockaddr
}

func (c *tcpConnection) Read(buf []byte) (int, error) {
	if _, err := syscall.Read(c.socket, buf); err != nil {
		return 0, err
	}

	idx := bytes.IndexByte(buf, 0)
	if idx == -1 {
		return len(buf), nil
	}
	return idx, io.EOF
}

func (c *tcpConnection) Write(p []byte) (int, error) {
	return syscall.Write(c.socket, p)
}

func (c *tcpConnection) Close() error {
	return syscall.Close(c.socket)
}
