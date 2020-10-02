package uuidgenerator

import (
	"crypto/rand"
	"fmt"
	"net/http"
)

// UUID is a custom multiplexer
type UUID struct {
}

// NewUUID return the reference to the UUID struct
func NewUUID() *UUID {
	return &UUID{}
}

func (p *UUID) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		giveRandomUUID(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func giveRandomUUID(w http.ResponseWriter, r *http.Request) {
	c := 10
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, fmt.Sprintf("%x", b))
}
