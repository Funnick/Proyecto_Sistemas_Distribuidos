package main

import (
	"log"
	"os"
	"server/server_side/server"
)

func main() {
	ip := os.Args[1]
	port := os.Args[2]
	port2 := os.Args[3]
	ipSucc := os.Args[4]
	name := os.Args[5]

	log.Println("Inicializando el servidor...")
	pl := server.NewPlatform(ip, port)
	pl.Run(port2, name, ipSucc)

}
