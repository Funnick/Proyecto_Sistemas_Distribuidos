package ap_client

import (
	"fmt"
	"github.com/urfave/cli"
	"strings"
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
	App.Version = "1.0.0"
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
			Value:       "./ap-client/config.cfg",
			Usage:       "Load configuration from `FILE`",
			Destination: &conf,
		},
	}
}

//Fix: Error que no permite cambiar el orden de -c FILE load-conf a load-conf -c FILE
func Commnads() {
	App.Commands = []cli.Command{
		{
			Name:    "get-all-agents",
			Aliases: []string{"A"},
			Usage:   "Request all agent to the server",
			Action: func(c *cli.Context) error {
				LoadConfig(conf)
				resp, agentList := GetAgentsRequest()
				fmt.Println(resp, agentList)
				return nil
			},
		},

		{
			Name:    "create-agent",
			Aliases: []string{"C"},
			Usage:   "Create new agent using the parameters [-i ip] [-pr port] [-des description] [-doc documentation]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "ip, i",
					Value:       "",
					Usage:       "Agent ip",
					Destination: &ip,
					Required:    true,
				},
				&cli.StringFlag{
					Name:        "port, pr",
					Value:       "",
					Usage:       "Agent port",
					Destination: &port,
					Required:    true,
				},
				&cli.StringFlag{
					Name:        "password, pass, p",
					Value:       "",
					Usage:       "Password that gives access to an agent",
					Destination: &password,
					Required:    true,
				},
				&cli.StringFlag{
					Name:        "description, des",
					Value:       "",
					Usage:       "Agent description",
					Destination: &description,
					Required:    true,
				},
				&cli.StringFlag{
					Name:        "doc",
					Value:       "",
					Usage:       "Agent doc",
					Destination: &doc,
					Required:    true,
				},
			},

			Action: func(c *cli.Context) error {
				LoadConfig(conf)
				err_val := create_validation()
				if err_val != nil {
					return err_val
				}
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
					Required:    true,
				},
				&cli.StringFlag{
					Name:        "password, pass, p",
					Value:       "",
					Usage:       "Password that gives access to an agent",
					Destination: &password,
					Required:    true,
				},
			},

			Action: func(c *cli.Context) error {
				fmt.Println(conf)
				LoadConfig(conf)

				err_val := delete_validation()
				if err_val != nil {
					return err_val
				}

				resp := DeleteAgentRequest(aid, password)
				fmt.Println(resp)

				return nil
			},
		},

		{
			Name:    "search-desc-agent",
			Aliases: []string{"Sd"},
			Usage:   "Search an Agent by description",
			Action: func(c *cli.Context) error {
				LoadConfig(conf)
				if c.NArg() > 0 {
					description := strings.Join(c.Args(), " ")
					resp, _ := SearchAgentDescRequest(description)
					fmt.Println(resp)
					return nil
				}
				fmt.Println("Please insert a valid string to search")
				return cli.NewExitError("Empty search string", 1)
			},
		},

		{
			Name:    "search-name-agent",
			Aliases: []string{"S"},
			Usage:   "Search an Agent by name",
			Action: func(c *cli.Context) error {
				LoadConfig(conf)
				if c.NArg() > 0 {
					description := strings.Join(c.Args(), " ")
					resp, _ := SearchAgentNameRequest(description)
					fmt.Println(resp)
					return nil
				}
				fmt.Println("Please insert a valid string to search")
				return cli.NewExitError("Empty search string", 1)
			},
		},

		{
			Name:    "update-agent",
			Aliases: []string{"U"},
			Usage:   "Update an existing agent using the parameters [-n name] [-p pswrd] [-i ip] [-pr port] [-np newpsword] [-des description] [-doc documentaion]",

			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "name, n",
					Value:       "",
					Usage:       "New Agent name",
					Destination: &name,
				},
				&cli.StringFlag{
					Name:        "password, pass, p",
					Value:       "",
					Usage:       "Password that gives access to an agent",
					Destination: &password,
				},
				&cli.StringFlag{
					Name:        "ip, i",
					Value:       "",
					Usage:       "New Agent ip",
					Destination: &ip,
				},
				&cli.StringFlag{
					Name:        "port, pr",
					Value:       "",
					Usage:       "New Agent port",
					Destination: &port,
				},
				&cli.StringFlag{
					Name:        "new-password, np",
					Value:       "",
					Usage:       "New Password that gives access to an agent",
					Destination: &newpassword,
				},
				&cli.StringFlag{
					Name:        "description, des",
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
				LoadConfig(conf)
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
