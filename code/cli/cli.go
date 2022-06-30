package main

import (
	client "agent_platform_client/ap-client"
	"fmt"
)

func main() {
	client.LoadConfig()
	_, l := client.GetAgentsRequest()
	r, k := client.SearchAgentRequest("000")
	fmt.Println(l)
	fmt.Println(r, k)

	//zz := ap_client.UpdateAgentRequest("Suma", "qwer")
}
