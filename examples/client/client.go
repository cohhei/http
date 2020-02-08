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
	// This method returns the response as io.Reader
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close() // The response should be closed.

	io.Copy(os.Stdout, resp)
}
