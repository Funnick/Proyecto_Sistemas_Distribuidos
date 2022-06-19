package main

import (
	"agent_platform_client/ap_client"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli"
)

var app = cli.NewApp()
var conf string
var aid string
var password string

func info() {
	app.Name = "Agent Platform Client CLI"
	app.Usage = "a CLI for the Agent Platform"
	app.Version = "0.0.0"
	app.Authors = []cli.Author{
		{
			Name:  "Miguel Alejandro Rodríguez Hernández",
			Email: "",
		},
		{
			Name:  "Manuel Antonio Vilas Valiente",
			Email: "",
		},
		{
			Name:  "Andrés Alejandro León Almaguer",
			Email: "",
		},
	}
}

func flags() {
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config, c",
			Value:       "",
			Usage:       "Load configuration from `FILE`",
			Destination: &conf,
		},
		&cli.StringFlag{
			Name:        "aid",
			Value:       "",
			Usage:       "Agent id",
			Destination: &aid,
		},
		&cli.StringFlag{
			Name:        "password, pass, p",
			Value:       "",
			Usage:       "Password that gives access to an agent",
			Destination: &password,
		},
	}
}

//Fix: Error que no permite cambiar el orden de -c FILE load-conf a load-conf -c FILE
func commnads() {
	app.Commands = []cli.Command{
		{
			Name:    "load-conf",
			Aliases: []string{"l"},
			Usage:   "Load a config txt",
			Action: func(c *cli.Context) error {
				if conf == "" {
					ap_client.LoadConfig()
				} else {
					//TODO: agregar una funcion que lea de un txt
					fmt.Println(conf)
				}
				return nil
			},
		},
        {
            Name: "get-agents",     
            Aliases: []string{"g"},
            Usage: "Request all agent to the server",
            Action: func (c *cli.Context) error {
                resp, agentList := ap_client.GetAgentsRequest()
                fmt.Println(resp, agentList) 
                return nil
            },
            
        },
        // {
        //     Name: "create-agent",     
        //     Aliases: []string{"c"},
        //     Usage: "Create new agent",
        //     Action: func (c *cli.Context) error {
        //         resp := ap_client.CreateAgentRequest() //TODO: ver como parsear de consola un CreateAgentMessage
        //         fmt.Println(resp) 
        //         return nil
        //     },
            
        // },
        {
            Name: "delete-agent",     
            Aliases: []string{"d"},
            Usage: "Delete an agent from the server",
            Action: func (c *cli.Context) error {
                resp := ap_client.DeleteAgentRequest(aid, password)  // TODO: ver como leer bien la contra y el id
                fmt.Println(resp) 
                return nil
            },
            
        },

        {
            Name: "search-agent",     
            Aliases: []string{"s"},
            Usage: "Search an Agent",
            Action: func (c *cli.Context) error {
                description := c.Args().Get(0)
                resp, _ := ap_client.SearchAgentRequest(description) //TODO: ver como parsear de consola un CreateAgentMessage
                fmt.Println(resp) 
                return nil
            },
            
        },
        // {
        //     Name: "update-agent",     
        //     Aliases: []string{"u"},
        //     Usage: "Update an existing agent",
        //     Action: func (c *cli.Context) error {
        //         resp := ap_client.UpdateAgentRequest() //TODO: ver como parsear de consola un CreateAgentMessage
        //         fmt.Println(resp) 
        //         return nil
        //     },
            
        // },
 
	}
}

func main() {
    
	info()
	flags()
	commnads()
    sort.Sort(cli.FlagsByName(app.Flags))
    sort.Sort(cli.CommandsByName(app.Commands))
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	} else {
		// ap_client.LoadConfig()
		// _, l := ap_client.GetAgentsRequest()
		// fmt.Println(l)
	}
}
