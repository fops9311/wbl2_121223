package main

import (
	_ "httprestapi/adapter"
	"httprestapi/server"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", server.Serve)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}
