# Too naive HTTP client

`cohhei/http` is a too simple and naive HTTP client that doesn't depend on `net/http`. This client has only one method, `Do` which sends an HTTP request.

## Usage

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

## Too naive TCP client

This package has a TCP client that depends on only `syscall`.

### TCP client usage

```go
// Create a client
addr := &http.Addr{
  IP:   [4]byte{127, 0, 0, 1},
  Port: 11111,
}
c := http.NewTCPClient()

// Connect to the address
if err := c.Connect(addr); err != nil {
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
addr := &Addr{
	IP:   [4]byte{8, 8, 8, 8},
	Port: 53,
}

b := []byte{0x01}
if err := c.SendTo(b, addr); err != nil {
	log.Fatal(err)
}
defer c.Close()

resp, err := c.Recvfrom()
if err != nil {
	log.Fatal(err)
}

io.Copy(os.Stdout, resp)
```