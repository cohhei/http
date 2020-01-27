package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
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
	Connect(addr *Addr) error
	Read() (io.Reader, error)
	Write(data []byte) (int, error)
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

// Connect creates a connection to the given IP address.
// If the client already has a connection, it will be closed.
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

// Read returns the data as `io.Reader`.
func (c *tcpclient) Read() (io.Reader, error) {
	buf := make([]byte, 1024)
	if _, err := syscall.Read(c.socket, buf); err != nil {
		log.Fatal(err)
	}
	return bytes.NewReader(buf), nil
}

// Write writes data.
func (c *tcpclient) Write(data []byte) (int, error) {
	return syscall.Write(c.socket, data)
}

// Close closes the connection.
func (c *tcpclient) Close() error {
	return syscall.Close(c.socket)
}
