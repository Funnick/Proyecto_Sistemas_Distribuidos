package main

import (
	"agent_platform_client/ap_client"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func main() {
	err := app.Run(os.Args)
	if err != nil {
        log.Fatal(err)
	} else {
		ap_client.LoadConfig()
		_, l := ap_client.GetAgentsRequest()
		fmt.Println(l)
	}
}
