package main

import (
	"io"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
)

func HelloWorld(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World, "+req.URL.Query().Get(":name")+"!\n")
}

func main() {
	m := pat.New()
	m.Get("/hello/:name", http.HandlerFunc(HelloWorld))

	http.Handle("/", m)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
