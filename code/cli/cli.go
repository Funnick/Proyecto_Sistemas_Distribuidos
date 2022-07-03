package main

import (
	client "client/ap_client"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
)

var app = client.App

func main() {
	//client.LoadConfig()
	client.Info()
	client.Flags()
	client.Commnads()
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	/*agent := client.CreateAgentMessage{Name: "Pepe", IP: "127.0.0.1", Port: "8080",
		Password: "123", Description: "Some", Documentation: "Documentation"}

	resp := client.CreateAgentRequest(agent)
	fmt.Println(resp)
	resp, agents := client.SearchAgentNameRequest("Pepe")
	fmt.Println(resp, agents)*/
}
