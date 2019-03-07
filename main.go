package main

import (
	"chatAppGo/trace"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	// "github.com/matryer/goblueprints/chapter1/trace"
)

// temp1 to represent a single template
type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

// ServeHTTP method handles the HTTP request .
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.temp1.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8081", "The addr of the application.")
	flag.Parse()

	// Setting up gomniauth
	gomniauth.SetSecurityKey("634743836443-be92hpdbpbfljrhln234j6a67j0phvge.apps.googleusercontent.com")
	gomniauth.WithProviders(
		google.New("634743836443-be92hpdbpbfljrhln234j6a67j0phvge.apps.googleusercontent.com", "9hnfpGjcVvLr2Om-lTSNM0fX",
			"http://localhost:8081/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	// If Bootstrap or other packages were being served with my own copy:
	// http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("/path/to/assets"))))

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()
	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
