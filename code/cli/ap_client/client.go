package ap_client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Request all agent to the server
func GetAgentsRequest() (resp string, agentList []Agent) {
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Get("http://127.0.0.1:5000/ap/agents")
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
	response, err := client.Post("http://127.0.0.1:5000/ap/create", "aplication/json", bytes.NewBuffer(messageJson))
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
	request, err := http.NewRequest(http.MethodDelete, "http://127.0.0.1:5000/ap/delete",
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
func SearchAgentRequest(description string) (resp string, agent Agent) {
	client := http.Client{Timeout: 10 * time.Second}
	var requestMessage SearchAgentMessage = SearchAgentMessage{
		Description: description,
	}
	jsonRequestMessage, err := json.Marshal(requestMessage)
	if err != nil {
		log.Fatal(err)
		return err.Error(), Agent{}
	}
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:5000/ap/search",
		bytes.NewBuffer(jsonRequestMessage))
	if err != nil {
		log.Fatal(err)
		return err.Error(), Agent{}
	}
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
		return err.Error(), Agent{}
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return err.Error(), Agent{}
	}
	var responseMessage SearchAgentMessageResponse
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		log.Fatal(err)
		return err.Error(), Agent{}
	}
	if len(responseMessage.Message) > 0 {
		return responseMessage.Message, Agent{}
	}
	return "", responseMessage.AgentFound
}

func UpdateAgentRequest(aid, password, newIP, newPort,
	newPassword, newDescription, newDocumentation string) (resp string) {
	client := http.Client{Timeout: 10 * time.Second}
	var requestMessage UpdateAgentMessage = UpdateAgentMessage{
		AID:              aid,
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
	request, err := http.NewRequest(http.MethodPut, "http://127.0.0.1:5000/ap/update",
		bytes.NewBuffer(jsonRequestMessage))
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
