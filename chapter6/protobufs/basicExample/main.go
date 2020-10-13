package main

import (
	"fmt"

	pb "github.com/lordkevinmo/hands-on-go/chapter6/protobufs/protofiles"
	"google.golang.org/protobuf/proto"
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

	p1 := &pb.Person{}
	body, _ := proto.Marshal(p)
	_ = proto.Unmarshal(body, p1)
	fmt.Println("Original struct loaded from proto file:", p)
	fmt.Println("Marshalled proto data: ", body)
	fmt.Println("Unmarshalled struct: ", p1)
}
