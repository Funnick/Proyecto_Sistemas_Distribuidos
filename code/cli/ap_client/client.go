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

func LoadConfig(path string) {

	readFile, err := os.Open(path)
	if err != nil {
		log.Println(err)
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
func GetAgentNamesRequest() (resp string, agentList []string) {
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(url + "/names")
	if err != nil {
		log.Println(err)
		return err.Error(), []string{}
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err.Error(), []string{}
	}
	var r GetAllResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Println(err)
		return err.Error(), []string{}
	}
	return "", r.ResponsesFound
}

func GetAgentDescsRequest() (resp string, agentList []string) {
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(url + "/descriptions")
	if err != nil {
		log.Println(err)
		return err.Error(), []string{}
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err.Error(), []string{}
	}
	var r GetAllResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Println(err)
		return err.Error(), []string{}
	}
	return "", r.ResponsesFound
}

// Create new agent
// TODO
func CreateAgentRequest(name, ip, port, password, description, documentation string) (resp string) {
	var agentMessage CreateAgentMessage = CreateAgentMessage{
		Name:          name,
		IP:            ip,
		Port:          port,
		Password:      password,
		Description:   description,
		Documentation: documentation,
	}
	client := http.Client{Timeout: 10 * time.Second}
	messageJson, err := json.Marshal(agentMessage)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	response, err := client.Post(url+"/create", "aplication/json", bytes.NewBuffer(messageJson))
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	var responseMessage ResponseMessage
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	return responseMessage.Message
}

// Delete an agent from the server
func DeleteAgentRequest(name string, password string) (resp string) {
	client := http.Client{Timeout: 10 * time.Second}
	var requestMessage DeleteAgentMessage = DeleteAgentMessage{
		Name: name, Password: password,
	}
	jsonRequestMessage, err := json.Marshal(requestMessage)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	request, err := http.NewRequest(http.MethodDelete, url+"/delete",
		bytes.NewBuffer(jsonRequestMessage))
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	var responseMessage ResponseMessage
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	return responseMessage.Message
}

// Search an Agent
func SearchAgentNameRequest(name string) (resp string, agent Agent) {
	client := http.Client{Timeout: 10 * time.Second}
	var requestMessage SearchAgentNameMessage = SearchAgentNameMessage{
		Name: name,
	}
	jsonRequestMessage, err := json.Marshal(requestMessage)
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}
	request, err := http.NewRequest(http.MethodGet, url+"/searchbyname",
		bytes.NewBuffer(jsonRequestMessage))
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}
	var responseMessage SearchAgentMessageResponse
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}
	if len(responseMessage.Message) > 0 {
		return responseMessage.Message, Agent{}
	}
	return "\u2713 OK", responseMessage.AgentsFound
}

func SearchAgentDescRequest(description string) (resp string, agent Agent) {
	client := http.Client{Timeout: 10 * time.Second}
	var requestMessage SearchAgentDescMessage = SearchAgentDescMessage{
		Description: description,
	}
	jsonRequestMessage, err := json.Marshal(requestMessage)
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}
	request, err := http.NewRequest(http.MethodGet, url+"/searchbydesc",
		bytes.NewBuffer(jsonRequestMessage))
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}
	var responseMessage SearchAgentMessageResponse
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}
	if len(responseMessage.Message) > 0 {
		return responseMessage.Message, Agent{}
	}
	return "\u2713 OK", responseMessage.AgentsFound
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
		log.Println(err)
		return
	}
	request, err := http.NewRequest(http.MethodPut, url+"/update",
		bytes.NewBuffer(jsonRequestMessage))
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	request.Header.Add("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	var responseMessage ResponseMessage
	err = json.Unmarshal(body, &responseMessage)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	return responseMessage.Message
}
