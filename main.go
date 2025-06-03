package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

//ServerHTTP handles the HTTP request

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, r)
}

func main() {
	// makes every randomly generated number unique
	rand.Seed(time.Now().UnixNano())
	var addr = flag.String("addr", ":8080", "Addr of the App")
	flag.Parse()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.Handle("/chat", &templateHandler{filename: "chat.html"})

	// Handle all websocket connections for chat rooms dynamically
	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		roomName := r.URL.Query().Get("room")
		if roomName == "" {
			http.Error(w, "Room name required", http.StatusBadRequest)
			return
		}
		realRoom := getRoom(roomName)
		realRoom.ServeHTTP(w, r)
	})

	log.Println("starting webserver on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
