package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	router := httprouter.New()
	// Mapping to paths
	router.ServeFiles("/static/*filepath", http.Dir("/Users/moiseagbenya/static"))
	log.Fatal(http.ListenAndServe(":8000", router))
}
