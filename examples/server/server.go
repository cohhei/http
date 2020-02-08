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
