package main

import (
	"log"

	"github.com/levigross/grequests"
)

func main() {
	resp, err := grequests.Get("http://httpbin.org/get", nil)
	// Request can be modified by passing an optional
	if err != nil {
		log.Fatalln("Unable to make request", err)
	}

	var returnData map[string]interface{}
	resp.JSON(&returnData)
	log.Println(returnData)
}
