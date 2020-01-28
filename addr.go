package http

import (
	"fmt"
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

func (a *Addr) sockaddr() syscall.Sockaddr {
	sa := &syscall.SockaddrInet4{
		Addr: a.IP,
		Port: a.Port,
	}
	return sa
}
