package main

import (
	"server/server_side/server"
)

func main() {
	server.StartServer("127.0.0.1", "5000")
}
