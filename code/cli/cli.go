package main

import (
	"agent_platform_client/ap_client"
	"fmt"
)

func main() {
	ap_client.LoadConfig()
	_, l := ap_client.GetAgentsRequest()
	r, k := ap_client.SearchAgentRequest("000")
	fmt.Println(l)
	fmt.Println(r, k)
}
