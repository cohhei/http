package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cohhei/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index.html"))
	})
	http.HandleFunc("/png", func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h["Content-Type"] = "image/png"

		f, err := os.Open("./examples/server/200x200.png")
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		defer f.Close()
		io.Copy(w, f)
	})
	http.HandleFunc("/ua", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%s", r.Header["User-Agent"])))
	})

	port := 8080
	log.Print(fmt.Sprintf("http://127.0.0.1:%d/", port))
	log.Print(fmt.Sprintf("http://127.0.0.1:%d/index.html", port))
	log.Print(fmt.Sprintf("http://127.0.0.1:%d/png", port))
	log.Print(fmt.Sprintf("http://127.0.0.1:%d/ua", port))
	if err := http.ListenAndServe(port); err != nil {
		log.Fatal(err)
	}
}
