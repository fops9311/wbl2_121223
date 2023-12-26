package server

import (
	"log"
	"net/http"
	"time"
)

func postJsonHandlerBuilder(f http.HandlerFunc) http.HandlerFunc {
	return logger(post(json(f)))
}
func getJsonHandlerBuilder(f http.HandlerFunc) http.HandlerFunc {
	return logger(get(json(f)))
}
func logger(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	}
}

//Набор утилит для организации http сервера
//Мне показался такой подход интересным, плюс он довольно быстрый согласно бенчмаркам

// get takes a HandlerFunc and wraps it to only allow the GET method
func get(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "GET")
}

// post takes a HandlerFunc and wraps it to only allow the POST method
func post(h http.HandlerFunc) http.HandlerFunc {
	return allowMethod(h, "POST")
}

// allowMethod takes a HandlerFunc and wraps it in a handler that only
//
// responds if the request method is the given method, otherwise it
//
// responds with HTTP 405 Method Not Allowed.
func allowMethod(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.Header().Set("Allow", method)
			http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h(w, r)
	}
}
func json(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h(w, r)
	}
}
