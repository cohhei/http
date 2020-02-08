package http

import (
	"io/ioutil"
	"testing"
)

func TestUDPClient(t *testing.T) {
	c := NewUDPClient()

	ip := [4]byte{8, 8, 8, 8}
	port := 53

	b := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01}
	if err := c.SendTo(b, ip, port); err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	r, err := c.Recvfrom()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp) == 0 {
		t.Fatal("response is empty")
	}
}
