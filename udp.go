package http

import (
	"io"
	"syscall"
)

// UDPClient is an interface for the UDP client
type UDPClient interface {
	// SendTo sends data.
	SendTo(data []byte, ip [4]byte, port int) error

	// Recvfrom returns io.Reader
	Recvfrom() (io.Reader, error)

	// Close closes the connection.
	Close() error
}

// NewUDPClient creates a client
func NewUDPClient() UDPClient {
	return &udpClient{}
}

type udpClient struct {
	socket int
}

func (c *udpClient) SendTo(data []byte, ip [4]byte, port int) error {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}
	c.socket = socket

	return syscall.Sendto(socket, data, 0, sockaddr(ip, port))
}

func (c *udpClient) Recvfrom() (io.Reader, error) {
	return &receiver{c.socket}, nil
}

func (c *udpClient) Close() error {
	return syscall.Close(c.socket)
}

type receiver struct {
	socket int
}

func (r *receiver) Read(buf []byte) (int, error) {
	n, _, err := syscall.Recvfrom(r.socket, buf, 0)
	if err != nil {
		return 0, err
	}

	return n, io.EOF
}
