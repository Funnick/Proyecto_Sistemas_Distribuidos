package main

import (
	"os"
	"server/chord"
)

func main() {
	ip := os.Args[1]
	port := os.Args[2]
	ipSucc := os.Args[3]
	portSucc := os.Args[4]
	name := os.Args[5]
	switch os.Args[6] {
	case "1":
		chord.JoinTest(ip, port, ipSucc, portSucc, name)
	case "2":
		chord.SleepTest(ip, port, ipSucc, portSucc, name)
	case "3":
		chord.StabilizeTest(ip, port, ipSucc, portSucc, name)
	case "4":
		chord.FireTest(ip, port, ipSucc, portSucc, name)
	default:
		chord.SelfTest(ip, port, ipSucc, portSucc, name)
	}
}
