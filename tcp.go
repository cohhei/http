package http

import (
	"bytes"
	"fmt"
	"io"
	"syscall"
)

// Addr contains an IP address and a port
type Addr struct {
	IP   [4]byte
	Port int
}

func (a *Addr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d:%d", a.IP[0], a.IP[1], a.IP[2], a.IP[3], a.Port)
}

func (a *Addr) sockaddrInet4() syscall.SockaddrInet4 {
	sa := syscall.SockaddrInet4{
		Addr: a.IP,
		Port: a.Port,
	}
	return sa
}

// TCPClient is an interface for the TPC client
type TCPClient interface {
	// Connect creates a TCP connection to the given IP address.
	// If the client already has a connection, it will be closed.
	// The tcp connection should be closed by `Close()`.
	Connect(addr *Addr) error

	// GetReader returns a reader to read responses
	GetReader() (io.Reader, error)

	// Send writes data.
	Send(data []byte) (int, error)

	// Close closes the connection.
	Close() error
}

// NewTCPClient creates a client
func NewTCPClient() TCPClient {
	return &tcpclient{}
}

type tcpclient struct {
	addr   *Addr
	socket int
}

func (c *tcpclient) Connect(addr *Addr) error {
	c.Close()

	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}
	c.socket = socket

	sa := addr.sockaddrInet4()
	if err := syscall.Connect(socket, &sa); err != nil {
		return err
	}
	return nil
}

func (c *tcpclient) GetReader() (io.Reader, error) {
	return &reader{c.socket}, nil
}

func (c *tcpclient) Send(data []byte) (int, error) {
	return syscall.Write(c.socket, data)
}

func (c *tcpclient) Close() error {
	return syscall.Close(c.socket)
}

type reader struct {
	socket int
}

func (r *reader) Read(buf []byte) (int, error) {
	if _, err := syscall.Read(r.socket, buf); err != nil {
		return 0, err
	}

	idx := bytes.IndexByte(buf, 0)
	if idx == -1 {
		return len(buf), nil
	}
	return idx + 1, io.EOF
}
