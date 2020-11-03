package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // Use default DB
	})

	pong, _ := client.Ping().Result() // Ignoring error
	fmt.Println(pong)
}
