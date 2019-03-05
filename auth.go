package main
import (
	"net/http"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// User was not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		// Handle any other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// In case of success, call the next handler
	h.next.ServeHTTP(w, r)
}

func MustAuth(hd http.Handler) http.Handler {
	return &authHandler{next: hd}
}