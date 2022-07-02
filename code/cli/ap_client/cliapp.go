package ap_client

import (
	"fmt"
	"log"
	"strings"

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
			Name:            "create-agent",
			Aliases:         []string{"C"},
			Usage:           "Create new agent",
			UsageText:       "!234",
			Description:     "una desc",
			ArgsUsage:       "kkk",
			SkipFlagParsing: false,
			HideHelp:        false,
			Hidden:          false,
			HelpName:        "doo!",
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

			Action: func(c *cli.Context) error {
                err_val := create_validation()
                if err_val  != nil {
                    return err_val
                }
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
                err_val := delete_validation()
                if err_val != nil{
                    return err_val
                }

				fmt.Println(password, aid)
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
				if c.NArg() > 0 {
					description := strings.Join(c.Args(), " ")
					log.Printf(description)
					resp, _ := SearchAgentRequest(description)
					fmt.Println(resp)
					return nil
				}
				fmt.Println("Please insert a valid string to search")
				return cli.NewExitError("Empty search string", 88)
			},
		},

		{
			Name:    "update-agent",
			Aliases: []string{"U"},
			Usage:   "Update an existing agent",

			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "name",
					Value:       "",
					Usage:       "New Agent name",
					Destination: &name,
				},
				&cli.StringFlag{
					Name:        "password",
					Value:       "",
					Usage:       "Password that gives access to an agent",
					Destination: &password,
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

			Action: func(c *cli.Context) error {
				fmt.Println(name, password, ip, port, newpassword, description, doc)
				err_valid := update_validation()
				if err_valid != nil {
					return err_valid
				}
				resp := UpdateAgentRequest(name, password, ip, port, newpassword, description, doc)
				fmt.Println(resp)
				return nil
			},
		},
	}
}

func update_validation() error {
	if name == "" {
		return cli.NewExitError("name error", 2)
	}
	if password == "" {
		return cli.NewExitError("Insert a correct Password", 2)
	}
	if ip == "" {
		return cli.NewExitError("Ip error", 2)
	}
	if port == "" {
		return cli.NewExitError("Port error", 2)
	}
	if newpassword == "" {
		return cli.NewExitError("New Password error", 2)
	}
	if description == "" {
		return cli.NewExitError("Description error", 2)
	}
	if doc == "" {
		return cli.NewExitError("Doc error", 2)
	}
	return nil
}

func create_validation() error {
	if ip == "" {
		return cli.NewExitError("Ip error", 2)
	}
	if port == "" {
		return cli.NewExitError("Port error", 2)
	}
	if password == "" {
		return cli.NewExitError("Password error", 2)
	}
	if description == "" {
		return cli.NewExitError("Description error", 2)
	}
	if doc == "" {
		return cli.NewExitError("Doc error", 2)
	}
	return nil
}

func delete_validation() error {
	if aid == "" {
		return cli.NewExitError("Id error", 2)
	}
	if password == "" {
		return cli.NewExitError("Password error", 2)
	}
	return nil
}
