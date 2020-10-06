package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

// Args is a type for arguments
type Args struct{}

// TimeServer is a reference to int64
type TimeServer int64

// GiveServerTime return a pointer to the reply with unix time
func (t *TimeServer) GiveServerTime(args *Args, reply *int64) error {
	// Fill reply pointer to send the data back
	*reply = time.Now().Unix()
	return nil
}

func main() {
	timeServer := new(TimeServer)
	rpc.Register(timeServer)
	rpc.HandleHTTP()
	// Listen for requests on port 1234
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
