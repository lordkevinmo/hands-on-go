package main

import (
	"encoding/json"
	"fmt"

	pb "github.com/lordkevinmo/hands-on-go/chapter6/protobufs/protofiles"
)

func main() {
	p := &pb.Person{
		ID:    1234,
		Name:  "Roger Federer",
		Email: "rf@email.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-4231", Type: pb.Person_HOME},
		},
	}

	body, _ := json.Marshal(p)
	fmt.Println(string(body))
}
