package main

import (
	"log"
	"net/http"

	controller "github.com/JeffNeff/webapp/backend/pkg"
)

func main() {
	c := &controller.Controller{}

	http.HandleFunc("/", c.RootHandler)
	http.HandleFunc("/event", c.HandlePost)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
