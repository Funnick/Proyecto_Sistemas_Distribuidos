package main

import (
	"dht/chord"
	"os"
)

func main() {
	ip := os.Args[1]
	port := os.Args[2]
	ipSucc := os.Args[3]
	portSucc := os.Args[4]
	if os.Args[5] == "1" {
		chord.JoinTest(ip, port, ipSucc, portSucc)
	}
	if os.Args[5] == "2" {
		chord.SleepTest(ip, port, ipSucc, portSucc)
	}
	if os.Args[5] == "3" {
		chord.StabilizeTest(ip, port, ipSucc, portSucc)
	}
	if os.Args[5] == "4" {
		chord.FireTest(ip, port, ipSucc, portSucc)
	} else {
		chord.SelfTest(ip, port, ipSucc, portSucc)
	}
}
