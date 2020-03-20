package http

import (
	"io/ioutil"
	"strings"
)

var (
	dnsServer = [4]byte{8, 8, 8, 8}
	dnsPort   = 53
	qHeader   = []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	qType     = []byte{0x00, 0x01}
	qClass    = []byte{0x00, 0x01}
)

func packet(domain string) []byte {
	domains := strings.Split(domain, ".")
	var r []byte
	r = append(r, qHeader...)
	for _, d := range domains {
		l := len(d)
		r = append(r, byte(l))
		r = append(r, d...)
	}
	r = append(r, 0x00)
	r = append(r, append(qType, qClass...)...)
	r = append(r, 0x00)
	return r
}

// DNSClient is an interface for the DNS client
type DNSClient interface {
	// Resolve returns an IP address from the given domain
	Resolve(string) ([4]byte, error)
}

// NewDNSClient creates a client
func NewDNSClient() DNSClient {
	return &dnsClient{udpClient: NewUDPClient()}
}

type dnsClient struct {
	udpClient UDPClient
}

func (c *dnsClient) Resolve(domain string) ([4]byte, error) {
	if domain == "127.0.0.1" || domain == "localhost" {
		return [4]byte{127, 0, 0, 1}, nil
	}

	p := packet(domain)
	if err := c.udpClient.SendTo(p, dnsServer, dnsPort); err != nil {
		return [4]byte{}, err
	}
	defer c.udpClient.Close()

	r, err := c.udpClient.Recvfrom()
	if err != nil {
		return [4]byte{}, err
	}

	resp, err := ioutil.ReadAll(r)
	if err != nil {
		return [4]byte{}, err
	}

	var ip [4]byte
	copy(ip[:], resp[len(resp)-4:])

	return ip, nil
}
