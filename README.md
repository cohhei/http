# Too naive TCP client

This is a TCP client that depends on only `syscall`.

## Usage

```go
// Create a client
addr := &Addr{
  IP:   [4]byte{127, 0, 0, 1},
  Port: 11111,
}
c := New()

// Connect to the address
if err := c.Connect(addr); err != nil {
  log.Fatal(err)
}

// Read data
data, err := c.Read()
if err != nil {
  log.Fatal(err)
}
b, err := ioutil.ReadAll(data)
if err != nil {
  log.Fatal(err)
}
fmt.Println(b)

// Write data
n, err := c.Write([]byte("world"))
if err != nil {
  log.Fatal(err)
}
```