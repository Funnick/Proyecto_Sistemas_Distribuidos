package main

import (
	client "agent_platform_client/ap-client"
	"log"
	"os"
	"sort"
	"github.com/urfave/cli"
)

var app = client.App

func main() {
	client.Info()
	client.Flags()
	client.Commnads()
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	} else {
		// ap_client.LoadConfig()
		// _, l := client.GetAgentsRequest()
		// fmt.Println(l)
	}
// func main() {
// 	client.LoadConfig()
// 	_, l := client.GetAgentsRequest()
// 	r, k := client.SearchAgentRequest("000")
// 	fmt.Println(l)
// 	fmt.Println(r, k)

// 	//zz := ap_client.UpdateAgentRequest("Suma", "qwer")
}
