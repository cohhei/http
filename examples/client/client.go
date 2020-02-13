package main

import (
	"fmt"
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
		Method: "GET",
		Host:   "example.com",
		Path:   "/index.html",
		Header: http.Header{},
	}

	// Send the request.
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
	io.Copy(os.Stdout, resp.Body)
}
