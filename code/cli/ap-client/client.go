package ap_client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var url string

func setURL(newURL string) {
	url = newURL
}

func LoadConfig() {
	readFile, err := os.Open("ap-client/config.cfg")
	if err != nil {
		log.Fatal(err)
		return
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	readFile.Close()
	for _, line := range lines {
		splitLine := strings.Split(line, " ")
		if splitLine[0] == "url" {
			setURL(splitLine[1])
		}
	}
}

// Request all agent to the server
func GetAgentsRequest() (resp string, agentList []Agent) {
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(url + "/agents")
	if err != nil {
		log.Fatal(err)
		return err.Error(), []Agent{}
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return err.Error(), []Agent{}
	}
	var agentArr []Agent
	err = json.Unmarshal(body, &agentArr)
	if err != nil {
		log.Fatal(err)
		return err.Error(), []Agent{}
	}
	return "", agentArr
}

// Create new agent
// TODO
func CreateAgentRequest(agentMessage CreateAgentMessage) (resp string) {
	client := http.Client{Timeout: 10 * time.Second}
	messageJson, err := json.Marshal(agentMessage)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	response, err := client.Post(url+"/create", "aplication/json", bytes.NewBuffer(messageJson))
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	var responseMessage ResponseMessage
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	return responseMessage.Message
}

// Delete an agent from the server
func DeleteAgentRequest(aid string, password string) (resp string) {
	client := http.Client{Timeout: 10 * time.Second}
	var requestMessage DeleteAgentMessage = DeleteAgentMessage{
		AID: aid, Password: password,
	}
	jsonRequestMessage, err := json.Marshal(requestMessage)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	request, err := http.NewRequest(http.MethodDelete, url+"/delete",
		bytes.NewBuffer(jsonRequestMessage))
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	var responseMessage ResponseMessage
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	return responseMessage.Message
}

// Search an Agent
func SearchAgentRequest(description string) (resp string, agent []Agent) {
	client := http.Client{Timeout: 10 * time.Second}
	var requestMessage SearchAgentMessage = SearchAgentMessage{
		Description: description,
	}
	jsonRequestMessage, err := json.Marshal(requestMessage)
	if err != nil {
		log.Fatal(err)
		return err.Error(), make([]Agent, 0)
	}
	request, err := http.NewRequest(http.MethodGet, url+"/search",
		bytes.NewBuffer(jsonRequestMessage))
	if err != nil {
		log.Fatal(err)
		return err.Error(), make([]Agent, 0)
	}
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return err.Error(), make([]Agent, 0)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return err.Error(), make([]Agent, 0)
	}
	var responseMessage SearchAgentMessageResponse
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		log.Fatal(err)
		return err.Error(), make([]Agent, 0)
	}
	if len(responseMessage.Message) > 0 {
		return responseMessage.Message, make([]Agent, 0)
	}
	return responseMessage.Message, responseMessage.AgentsFound
}

func UpdateAgentRequest(name, password, newIP, newPort,
	newPassword, newDescription, newDocumentation string) (resp string) {
	client := http.Client{Timeout: 10 * time.Second}
	var requestMessage UpdateAgentMessage = UpdateAgentMessage{
		Name:             name,
		Password:         password,
		NewIP:            newIP,
		NewPort:          newPort,
		NewPassword:      newPassword,
		NewDescription:   newDescription,
		NewDocumentation: newDocumentation,
	}
	jsonRequestMessage, err := json.Marshal(requestMessage)
	if err != nil {
		log.Fatal(err)
		return
	}
	request, err := http.NewRequest(http.MethodPut, url+"/update",
		bytes.NewBuffer(jsonRequestMessage))
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	var responseMessage ResponseMessage
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		log.Fatal(err)
		return err.Error()
	}
	return responseMessage.Message
}

func updARIP()            {}
func updARPort()          {}
func updARName()          {}
func updARDescription()   {}
func updARDocumentation() {}
