package ap_client

import (
	"fmt"
	"github.com/urfave/cli"

)

var conf string
var aid string
var password string
var ip string
var port string
var description string
var doc string
var name string
var newpassword string

var App = cli.NewApp()

func Info() {
	App.Name = "Agent Platform Client CLI"
	App.Usage = "a CLI for the Agent Platform"
	App.Version = "0.0.0"
	App.Authors = []cli.Author{
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

func Flags() {
	App.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "config, c",
			FilePath:    "./ap_client/config.txt",
			Usage:       "Load configuration from `FILE`",
			Destination: &conf,
		},
	}
}

//Fix: Error que no permite cambiar el orden de -c FILE load-conf a load-conf -c FILE
func Commnads() {
	App.Commands = []cli.Command{
		{
			Name:    "load-conf",
			Aliases: []string{"L"},
			Usage:   "Load a config txt",
			Action: func(c *cli.Context) error {
				conf_path := c.Args().First()
				if conf_path == "" {
					LoadConfig()
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
				resp, agentList := GetAgentsRequest()
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
		        resp := CreateAgentRequest(ip, port, password, description, doc) 
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
				resp := DeleteAgentRequest(aid, password)
				fmt.Println(resp)
                
				return nil
			},
		},

		{
			Name:    "search-agent",
			Aliases: []string{"S"},
			Usage:   "Search an Agent",
			Action: func(c *cli.Context) error {
                description := c.Args().Get(0) //TODO: agregar un flag booleano o un StringSliceFlag para parsear la lista o permitir recibir mas de un string
				resp, _ := SearchAgentRequest(description)
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
		        resp := UpdateAgentRequest(name, password, ip, port, newpassword, description, doc) 
		        fmt.Println(resp)
		        return nil
		    },

		},

	}
}
