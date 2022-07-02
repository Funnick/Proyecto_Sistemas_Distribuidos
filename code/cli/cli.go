package main

import (
	"fmt"
	client "server/cli/ap_client"
)

func main() {
	client.LoadConfig()
	agent := client.CreateAgentMessage{Name: "Pepe", IP: "127.0.0.1", Port: "8080",
		Password: "123", Description: "Some", Documentation: "Documentation"}
	resp := client.CreateAgentRequest(agent)
	fmt.Println(resp)
	resp, agents := client.SearchAgentNameRequest("Pepe")
	fmt.Println(resp, agents)

	//resp = client.DeleteAgentRequest("Pepe", "123")
	//fmt.Println(resp)
	//_, l := client.GetAgentsRequest()
	//r, k := client.SearchAgentRequest("000")
	//fmt.Println(l)
	//fmt.Println(r, k)

	//zz := ap_client.UpdateAgentRequest("Suma", "qwer")
}
