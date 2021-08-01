package main

import (
	"embed"
	"fmt"
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

func WebComponent(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	username := ""
	for _, cookie := range cookies {
		if cookie.Name == "username" {
			username = cookie.Value
			break
		}
	}
	if username == "" {
		username = "guest"
	}
	w.Header().Add("Content-Type", "text/javascript; charset=UTF-8")
	w.Write([]byte(fmt.Sprintf(`
		class LoggedIn extends HTMLElement {
			constructor() {
				super()
				const username = "%s"
				const shadow = this.attachShadow({ mode: 'open' });
				const container = document.createElement('div');
				container.innerHTML = "<h1>Hello </h1>" +
					"<h2> Name: " + username + "</h2>" +
					"<h2> Reversed: " + username.split("").reverse().join("") + "</h2>";
				shadow.appendChild(container)
			}
		}
		customElements.define('logged-in', LoggedIn)
		`, username)))
}

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.FS(indexHTML)))
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/logout", Logout).Methods("POST")
	r.HandleFunc("/webcomponent.js", WebComponent).Methods("GET")

	log.Println("server starting on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
