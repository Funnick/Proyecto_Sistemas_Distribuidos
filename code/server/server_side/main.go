package main

import (
	"server/server_side/server"
)

func main() {
	pl := server.NewPlatform("127.0.0.1", "6000")
	pl.Run("6001")
}
