package main

import (
	"fmt"
	"io"
	"log"

	"github.com/cohhei/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/index.html", func(w io.Writer, r *http.Request) {
		fmt.Println(r)
		w.Write([]byte("HTTP/1.1 200 OK\naaaaa: bbbbbb\n\nHello World!"))
		fmt.Println("done")
	})
	log.Print("http://127.0.0.1:8080/index.html")
	if err := http.ListenAndServe(8080); err != nil {
		log.Fatal(err)
	}
}
