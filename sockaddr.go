package http

import "syscall"

func sockaddr(ip [4]byte, port int) syscall.Sockaddr {
	sa := &syscall.SockaddrInet4{
		Addr: ip,
		Port: port,
	}
	return sa
}
