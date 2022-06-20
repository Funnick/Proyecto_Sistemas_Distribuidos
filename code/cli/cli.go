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
			FilePath:    "./ap_client/config.txt",
			Usage:       "Load configuration from `FILE`",
			Destination: &conf,
		},
	}
}

//Fix: Error que no permite cambiar el orden de -c FILE load-conf a load-conf -c FILE
func commnads() {
	app.Commands = []cli.Command{
		{
			Name:    "load-conf",
			Aliases: []string{"L"},
			Usage:   "Load a config txt",
			Action: func(c *cli.Context) error {
				conf_path := c.Args().First()
				if conf_path == "" {
					ap_client.LoadConfig()
				} else {
					//TODO: agregar una funcion que lea de un txt
					fmt.Println(conf)
				}
				return nil
			},
		},
		{
			Name:    "get-agents",
			Aliases: []string{"G"},
			Usage:   "Request all agent to the server",
			Action: func(c *cli.Context) error {
				resp, agentList := ap_client.GetAgentsRequest()
				fmt.Println(resp, agentList)
				return nil
			},
		},
		// {
		//     Name: "create-agent",
		//     Aliases: []string{"C"},
		//     Usage: "Create new agent",
		//     Action: func (c *cli.Context) error {
		//         resp := ap_client.CreateAgentRequest() //TODO: ver como parsear de consola un CreateAgentMessage
		//         fmt.Println(resp)
		//         return nil
		//     },

		// },
		{
			Name:    "delete-agent",
			Aliases: []string{"D"},
			Usage:   "Delete an agent from the server using [-a aid] [-p pswrd] parameters",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "aid, a",
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
			},

			Action: func(c *cli.Context) error {
                fmt.Println(password , aid)
				resp := ap_client.DeleteAgentRequest(aid, password) // TODO: ver como leer bien la contra y el id
				fmt.Println(resp)
                
				return nil
			},
		},

		{
			Name:    "search-agent",
			Aliases: []string{"S"},
			Usage:   "Search an Agent",
			Action: func(c *cli.Context) error {
                description := c.Args().Get(0) //TODO: agregar un flag booleano para parsear la lista o permitir recibir mas de un string
				resp, _ := ap_client.SearchAgentRequest(description) //TODO: ver como parsear de consola un CreateAgentMessage
				fmt.Println(resp)
				return nil
			},
		},
		// {
		//     Name: "update-agent",
		//     Aliases: []string{"U"},
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
