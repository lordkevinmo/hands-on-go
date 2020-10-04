package main

import (
	"log"
	"net/rpc"
)

// Args is a arguments type definition
type Args struct{}

func main() {
	var reply int64
	args := Args{}
	client, err := rpc.DialHTTP("tcp", "localhost"+":1234")
	if err != nil {
		log.Println("dialing:", err)
	}
	err = client.Call("TimeServer.GiveServerTime", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	log.Printf("%d", reply)
}
