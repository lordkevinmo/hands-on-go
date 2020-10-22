package main

import (
	"flag"
	"log"
)

var fname = flag.String("name", "stranger", "your wonderful name")
var age = flag.Int("age", 0, "Your graceful age")

func main() {
	flag.Parse()
	log.Printf("Hello %s (%d years), Welcome to the command line world", *fname, *age)
}
