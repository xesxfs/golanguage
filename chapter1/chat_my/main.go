package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates_my", t.filename)))

	})
	t.templ.Execute(w, r)

}

func main() {
	var addr = flag.String("addr", ":8080", "The address of the application")
	flag.Parse()
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte(`
	//        <html>
	//          <head>
	//             <title>Chat</title>
	//          </head>
	//          <body>
	//          Let' chat!
	//          </body>
	//        </html>

	//        `))
	// })

	r := newRoom()

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer:", err)

	}

}
