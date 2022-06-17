package main

import (
	"agent_platform_client/ap_client"
	"fmt"
)

func main() {
	ap_client.LoadConfig()
	_, l := ap_client.GetAgentsRequest()
	fmt.Println(l)
}
