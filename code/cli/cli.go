package main

import (
	client "agent_platform_client/ap-client"
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
var ip string
var port string
var description string
var doc string
var name string
var newpassword string

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
					client.LoadConfig()
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
				resp, agentList := client.GetAgentsRequest()
				fmt.Println(resp, agentList)
				return nil
			},
		},

		{
		    Name: "create-agent",
		    Aliases: []string{"C"},
		    Usage: "Create new agent",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "ip",
					Value:       "",
					Usage:       "Agent ip",
					Destination: &ip,
				},
				&cli.StringFlag{
					Name:        "port",
					Value:       "",
					Usage:       "Agent port",
					Destination: &port,
				},
				&cli.StringFlag{
					Name:        "password, pass, p",
					Value:       "",
					Usage:       "Password that gives access to an agent",
					Destination: &password,
				},
				&cli.StringFlag{
					Name:        "description",
					Value:       "",
					Usage:       "Agent description",
					Destination: &description,
				},
				&cli.StringFlag{
					Name:        "doc",
					Value:       "",
					Usage:       "Agent doc",
					Destination: &doc,
				},
			},

		    Action: func (c *cli.Context) error {
                fmt.Println(ip, port, password, description, doc)
		        resp := client.CreateAgentRequest(ip, port, password, description, doc) 
		        fmt.Println(resp)
		        return nil
		    },

		},
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
				resp := client.DeleteAgentRequest(aid, password) // TODO: ver como leer bien la contra y el id
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
				resp, _ := client.SearchAgentRequest(description) //TODO: ver como parsear de consola un CreateAgentMessage
				fmt.Println(resp)
				return nil
			},
		},

		{
		    Name: "update-agent",
		    Aliases: []string{"U"},
		    Usage: "Update an existing agent",

			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "name",
					Value:       "",
					Usage:       "New Agent name",
					Destination: &name,
				},
				&cli.StringFlag{
					Name:        "ip",
					Value:       "",
					Usage:       "New Agent ip",
					Destination: &ip,
				},
				&cli.StringFlag{
					Name:        "port",
					Value:       "",
					Usage:       "New Agent port",
					Destination: &port,
				},
				&cli.StringFlag{
					Name:        "new-password",
					Value:       "",
					Usage:       "New Password that gives access to an agent",
					Destination: &newpassword,
				},
				&cli.StringFlag{
					Name:        "description",
					Value:       "",
					Usage:       "New Agent description",
					Destination: &description,
				},
				&cli.StringFlag{
					Name:        "doc",
					Value:       "",
					Usage:       "New Agent doc",
					Destination: &doc,
				},
			},

		    Action: func (c *cli.Context) error {
                fmt.Println(name, password, ip, port, newpassword, description, doc)
		        resp := client.UpdateAgentRequest(name, password, ip, port, newpassword, description, doc) 
		        fmt.Println(resp)
		        return nil
		    },

		},

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
