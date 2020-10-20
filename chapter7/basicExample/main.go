package main

import (
	"log"

	"github.com/lordkevinmo/hands-on-go/chapter7/basicExample/helper"
)

func main() {
	_, err := helper.InitDB()
	if err != nil {
		log.Println(err)
	}

	log.Println("Database tables are successfully initiated!")
}
