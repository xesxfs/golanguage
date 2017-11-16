package main

import (
	_ "chapter2/trace_my"
	"flag"
	"github.com/stretchr/gomniauth"
	_ "github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	_ "github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"log"
	"net/http"
	_ "os"
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
	data := map[string]interface{}{"Host": r.Host}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)

	}
	t.templ.Execute(w, data)

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

	gomniauth.SetSecurityKey("98dfbg7iu2nb4uywevihjw4tuiyub34noilk")
	gomniauth.WithProviders(
		// facebook.New("key", "secret", "http://localhost:8080/auth/callback/facebook"),
		github.New("b1bf4cb03dc2afa37611", "92771e2ce33888ae07f6c22d6c3d903777f7f116",
			"http://localhost:8080/auth/callback/github"),
		// google.New("key", "secret", "http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	// r.tracer = trace_my.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServer:", err)

	}

}
