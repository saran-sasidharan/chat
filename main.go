package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	port := flag.String("port", "", "Port to listen to")
	flag.Parse()

	if *port == "" {
		p := os.Getenv("PORT")
		if p == "" {
			log.Fatal("port/host not provided")
		}
		*port = p
	}

	r := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	go r.run()

	log.Println("Starting web server on", *port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
