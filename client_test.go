package tcp

import (
	"bytes"
	"io/ioutil"
	"syscall"
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	ch := make(chan struct{})
	go func() {
		socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
		if err != nil {
			t.Fatal(err)
		}
		defer syscall.Close(socket)

		sa := syscall.SockaddrInet4{
			Addr: [4]byte{127, 0, 0, 1},
			Port: 11111,
		}
		if err := syscall.Bind(socket, &sa); err != nil {
			t.Fatal(err)
		}
		if err := syscall.Listen(socket, 5); err != nil {
			t.Fatal(err)
		}

		socket1, _, err := syscall.Accept(socket)
		if err != nil {
			t.Fatal(err)
		}
		defer syscall.Close(socket1)

		n, err := syscall.Write(socket1, []byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
		if n != 5 {
			t.Fatal(n)
		}

		buf := make([]byte, 5)
		if _, err := syscall.Read(socket1, buf); err != nil {
			t.Fatal(err)
		}
		expected := []byte("world")
		if bytes.Compare(buf, expected) != 0 {
			t.Errorf("expected: '%v', actual: '%v'\n", expected, buf)
		}

		ch <- struct{}{}
	}()

	time.Sleep(10 * time.Millisecond)
	addr := &Addr{
		IP:   [4]byte{127, 0, 0, 1},
		Port: 11111,
	}
	c := New()
	if err := c.Connect(addr); err != nil {
		t.Fatal(err)
	}
	data, err := c.Read()
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(data)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte("helloa")
	if bytes.Compare(b, expected) == 0 {
		t.Errorf("expected: '%v', actual: '%v'\n", expected, b)
	}
	n, err := c.Write([]byte("world"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Fatal(n)
	}
	<-ch
}
