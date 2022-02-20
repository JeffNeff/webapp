package main

import (
	"net/http"

	controller "github.com/JeffNeff/webapp/backend/pkg"
)

func main() {

	c := &controller.Controller{}

	http.HandleFunc("/", c.RootHandler)

	http.ListenAndServe(":8080", nil)
}
