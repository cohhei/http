package http

import "testing"

func TestPacket(t *testing.T) {
	got := packet("example.com")
	want := []byte{0x00, 0x00, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00}
	if string(got) != string(want) {
		t.Fatalf("\nwant:\t%v\ngot:\t%v\n", want, got)
	}
}

func TestResolve(t *testing.T) {
	c := NewDNSClient()
	got, err := c.Resolve("example.com")
	if err != nil {
		t.Fatal(err)
	}
	want := [4]byte{93, 184, 216, 34}
	if got != want {
		t.Fatalf("\nwant:\t%v\ngot:\t%v\n", want, got)
	}
}
