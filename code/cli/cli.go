package main

import (
	"agent_platform_client/ap_client"
	"fmt"
)

func main() {
	/*
		agentMessage := &ap_client.CreateAgentMessage{IP: "10.6.0.1", Port: "8050",
			Password: "myPassword", Description: "myDescription", Documentation: "myDoc"}
		s := ap_client.CreateAgentRequest(*agentMessage)
		fmt.Println(s)
	*/
	mess, ag := ap_client.GetAgentsRequest()
	if len(mess) > 0 {
		fmt.Println(mess)
	}
	fmt.Println(ag)
}
