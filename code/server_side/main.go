package main

import (
	"agent_platform/server"
)

func main() {
	server.StartServer("127.0.0.1", "5000")
}
