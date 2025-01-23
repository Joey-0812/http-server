package main

import (
	"log"
	"net/http"
)

func handlerHttps(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, HTTPS with self-generated certificate!"))
}

func handlerHttp(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, HTTP request!"))
}

func main() {
	go func() {
		http.HandleFunc("/http", handlerHttp)
		log.Fatal(http.ListenAndServe(":80", nil))
	}()
	http.HandleFunc("/https", handlerHttps)
	log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", nil))
}
