package main

import (
	"chatAppGo/trace"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

var avatars Avatar = TryAvatars{
	UseFileSystemAvatar,
	UseAuthAvatar,
	UseGravatar}

// temp1 to represent a single template
type templateHandler struct {
	filename string
	templ    *template.Template
}

// ServeHTTP method handles the HTTP request .
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if t.templ == nil {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	}

	data := map[string]interface{}{
		"Host": r.Host,
	}

	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

var host = flag.String("host", ":8081", "The host of the application.")

func main() {
	flag.Parse()

	gomniauth.SetSecurityKey("CKzo4DBYZXgLWC4Y7qREY6losI1jgOykItKqlH3bezfM5G9cxZk2EaSm9LJaJ0dC")
	gomniauth.WithProviders(google.New("634743836443-be92hpdbpbfljrhln234j6a67j0phvge.apps.googleusercontent.com", "9hnfpGjcVvLr2Om-lTSNM0fX", "http://localhost:8081/auth/callback/google"))

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)

	http.Handle("/avatars/", http.StripPrefix("/avatars", http.FileServer(http.Dir("./avatars"))))

	go r.run()
	// start the web server
	log.Println("Starting web server on", *host)
	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
