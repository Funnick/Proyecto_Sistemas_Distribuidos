package main

import (
	"server/server"
)

func main() {
	pl := server.NewPlatform("127.0.0.1", "5000")
	pl.Run()
}
