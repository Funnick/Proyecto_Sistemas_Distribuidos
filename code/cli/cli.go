package main

import (
	client "client/ap_client"
	"log"
	"os"
)

var app = client.App

func main() {
	client.AppInit()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
