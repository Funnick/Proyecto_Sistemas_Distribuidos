package main

import (
	"fmt"
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
	go pl.Run(port2, name, ipSucc)
	var command string
	for {
		fmt.Scanln(&command)
		switch command {
		case "stop":
			pl.Stop()
		case "join":
			var ip string
			var port string
			fmt.Print("IP: ")
			fmt.Scanln(&ip)
			fmt.Print("Port: ")
			fmt.Scanln(&port)
			pl.Join(ip, port)
		default:
			continue
		}
	}

	/*channel := make(chan string)


	go func(channel chan string) {
		for {
			var command string
			fmt.Scanln(&command)
			channel <- name
		}
	}(channel)
	*/
}
