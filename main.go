package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed index.html
var indexHTML embed.FS

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error"))
		return
	}
	username := r.FormValue("username")

	http.SetCookie(w, &http.Cookie{
		Name:  "username",
		Value: username,
	})
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Value:  "",
		MaxAge: -1,
	})
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.FS(indexHTML)))
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("POST")

	log.Println("server starting on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
