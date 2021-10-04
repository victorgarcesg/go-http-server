package main

import (
	"crypto/subtle"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	CONN_HOST      = "localhost"
	CONN_PORT      = "8080"
	ADMIN_USER     = "ADMIN"
	ADMIN_PASSWORD = "admin"
)

func BasicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		realm := "Please enter your username and password"
		user, pass, ok := r.BasicAuth()
		if !ok ||
			subtle.ConstantTimeCompare([]byte(user), []byte(ADMIN_USER)) != 1 ||
			subtle.ConstantTimeCompare([]byte(pass), []byte(ADMIN_PASSWORD)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm=`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("You are Unauthorized to access the application. \n"))
			return
		}
		handler(w, r)
	}
}

var GetRequestHandler = http.HandlerFunc(
	func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World!"))
	})

var PostRequestHandler = http.HandlerFunc(
	func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("It's a Post Request!"))
	})

var PathVariableHandler = http.HandlerFunc(
	func(rw http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		rw.Write([]byte("Hi " + name))
	})

func main() {
	router := mux.NewRouter()
	router.Handle("/", handlers.LoggingHandler(os.Stdout, BasicAuth(GetRequestHandler))).Methods("GET")
	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
	router.Handle("/post", handlers.LoggingHandler(logFile, BasicAuth(PostRequestHandler))).Methods("POST")
	router.Handle("/hello/{name}", handlers.CombinedLoggingHandler(logFile, BasicAuth(PathVariableHandler))).Methods("GET", "PUT")
	http.ListenAndServe(CONN_HOST+":"+CONN_PORT, router)
}
