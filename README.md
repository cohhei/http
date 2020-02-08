# Too naive HTTP client and server

`cohhei/http` is a too simple and naive HTTP client and server library that doesn't depend on `net/http`. This client has only one method, `Do` which sends an HTTP request.

## Usage

### HTTP client

```go
package main

import (
	"io"
	"log"
	"os"

	"github.com/cohhei/http"
)

func main() {
	// Create a client.
	c := http.NewClient()

	// Create a request.
	req := &http.Request{
		Method:  "GET",
		Host:    "example.com",
		Path:    "/index.html",
		Headers: map[string]string{},
		Body:    "",
	}

	// Send the request.
	// This method returns the response as io.ReadCloser
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close() // The response should be closed.

	io.Copy(os.Stdout, resp)
}

```

### HTTP server

```go
package main

import (
	"io"
	"log"
	"os"

	"github.com/cohhei/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/", func(w io.Writer, r io.Reader) {
		io.Copy(os.Stderr, r)
		w.Write([]byte("HTTP/1.1 200 OK\naaaaa: bbbbbb\n\nHello World!"))
	})
	log.Print("http://127.0.0.1:8080/")
	if err := http.ListenAndServe(8080); err != nil {
		log.Fatal(err)
	}
}
```

## Too naive TCP client

This package has a TCP client that depends on only `syscall`.

### TCP client usage

```go
// Create a client
c := http.NewTCPClient()

// Connect to the address
ip := [4]byte{127, 0, 0, 1}
port := 11111
if err := c.Connect(ip, port); err != nil {
  	log.Fatal(err)
}

// Read data
resp, err := c.Read()
if err != nil {
	log.Fatal(err)
}
io.Copy(os.Stdout, resp)

// Write data
n, err := c.Write([]byte("world"))
if err != nil {
	log.Fatal(err)
}
```

## Too naive UDP client

This package has a UDP client that depends on only `syscall`.

### UDP client usage

```go
c := http.NewUDPClient()
ip := [4]byte{8, 8, 8, 8}
port := 53
b := []byte{0x01}
if err := c.SendTo(b, ip, port); err != nil {
	log.Fatal(err)
}
defer c.Close()

resp, err := c.Recvfrom()
if err != nil {
	log.Fatal(err)
}

io.Copy(os.Stdout, resp)
```