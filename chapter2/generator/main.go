package main

import (
	"net/http"

	"github.com/lordkevinmo/hands-on-go/chapter2/uuidgenerator"
)

func main() {
	mux := uuidgenerator.NewUUID()
	http.ListenAndServe(":8080", mux)
}
