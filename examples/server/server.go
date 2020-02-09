package main

import (
	"fmt"
	"io"
	"log"

	"github.com/cohhei/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/", func(w io.Writer, r *http.Request) {
		w.Write([]byte("HTTP/1.1 200 OK\n\nHello World!"))
	})
	http.HandleFunc("/index.html", func(w io.Writer, r *http.Request) {
		w.Write([]byte("HTTP/1.1 200 OK\n\nindex.html"))
	})
	http.HandleFunc("/ua", func(w io.Writer, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\n\n%s", r.Headers["User-Agent"])))
	})
	log.Print("http://127.0.0.1:8080/")
	log.Print("http://127.0.0.1:8080/index.html")
	log.Print("http://127.0.0.1:8080/ua")
	if err := http.ListenAndServe(8080); err != nil {
		log.Fatal(err)
	}
}
