package main

import (
	"Agent_Platform/server"
	"fmt"
)


func main(){
	addr := server.Address { IP: "127.0.0.1", Port:"8080" }
	fmt.Println(addr.IP)
}