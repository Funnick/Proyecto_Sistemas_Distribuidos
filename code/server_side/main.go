package main

import (
	"agent_platform/server"
)

func main() {
	pl := server.NewPlatform("127.0.0.1", "5000")
	pl.Run()
}
