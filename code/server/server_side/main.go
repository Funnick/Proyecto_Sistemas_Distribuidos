package main

import (
	"os"
	"server/server_side/server"
)

func main() {
	ip := os.Args[1]
	port := os.Args[2]
	port2 := os.Args[3]
	ipSucc := os.Args[4]
	portSucc := os.Args[5]
	name := os.Args[6]

	pl := server.NewPlatform(ip, port)
	pl.Run(port2, name, ipSucc, portSucc)

}
