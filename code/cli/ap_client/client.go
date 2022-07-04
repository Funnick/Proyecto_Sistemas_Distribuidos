package ap_client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var url []string

func setURL(newURL []string) {
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
			url = append(url, splitLine[1])
		}
	}
}

// Request all agent to the server
func GetAgentNamesRequest() (resp string, agentList []string) {
	client := http.Client{Timeout: 5 * time.Second}
	count := len(url)
	var _err error
	var httpResp *http.Response
	for i := 0; i < len(url); i++ {
		response, err := client.Get(url[i] + "/names")
		if err != nil {
			count--
			_err = err
			continue
		} else {
			httpResp = response
			break
		}
	}
	if count == 0 {
		log.Println(_err.Error())
		return _err.Error(), []string{}
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
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
	client := http.Client{Timeout: 5 * time.Second}
	count := len(url)
	var _err error
	var httpResp *http.Response
	for i := 0; i < len(url); i++ {
		response, err := client.Get(url[i] + "/descriptions")
		if err != nil {
			count--
			_err = err
			continue
		} else {
			httpResp = response
			break
		}
	}
	if count == 0 {
		log.Println(_err.Error())
		return _err.Error(), []string{}
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
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
	client := http.Client{Timeout: 5 * time.Second}
	messageJson, err := json.Marshal(agentMessage)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	count := len(url)
	var _err error
	var httpResp *http.Response
	for i := 0; i < len(url); i++ {
		response, err := client.Post(url[i]+"/create", "aplication/json", bytes.NewBuffer(messageJson))
		if err != nil {
			count--
			_err = err
			continue
		} else {
			httpResp = response
			break
		}
	}
	if count == 0 {
		log.Println(_err.Error())
		return _err.Error()
	}

	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
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
	client := http.Client{Timeout: 5 * time.Second}
	var requestMessage DeleteAgentMessage = DeleteAgentMessage{
		Name: name, Password: password,
	}
	jsonRequestMessage, err := json.Marshal(requestMessage)
	if err != nil {
		log.Println(err)
		return err.Error()
	}

	count := len(url)
	var _err error
	var httpResp *http.Response
	for i := 0; i < len(url); i++ {
		request, _ := http.NewRequest(http.MethodDelete, url[i]+"/delete",
			bytes.NewBuffer(jsonRequestMessage))
		request.Header.Add("Accept", "application/json")
		response, err := client.Do(request)
		if err != nil {
			count--
			_err = err
		} else {
			httpResp = response
			break
		}
	}
	if count == 0 {
		log.Println(_err.Error())
		return _err.Error()
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
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
	client := http.Client{Timeout: 5 * time.Second}
	var requestMessage SearchAgentNameMessage = SearchAgentNameMessage{
		Name: name,
	}
	jsonRequestMessage, err := json.Marshal(requestMessage)
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}

	count := len(url)
	var _err error
	var httpResp *http.Response
	for i := 0; i < len(url); i++ {
		request, _ := http.NewRequest(http.MethodGet, url[i]+"/searchbyname",
			bytes.NewBuffer(jsonRequestMessage))
		request.Header.Add("Accept", "application/json")
		response, err := client.Do(request)
		if err != nil {
			count--
			_err = err
			continue
		} else {
			httpResp = response
			break
		}
	}
	if count == 0 {
		log.Println(_err.Error())
		return _err.Error(), Agent{}
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
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
	client := http.Client{Timeout: 5 * time.Second}
	var requestMessage SearchAgentDescMessage = SearchAgentDescMessage{
		Description: description,
	}
	jsonRequestMessage, err := json.Marshal(requestMessage)
	if err != nil {
		log.Println(err)
		return err.Error(), Agent{}
	}

	count := len(url)
	var _err error
	var httpResp *http.Response
	for i := 0; i < len(url); i++ {
		request, _ := http.NewRequest(http.MethodGet, url[i]+"/searchbydesc",
			bytes.NewBuffer(jsonRequestMessage))
		request.Header.Add("Accept", "application/json")
		response, err := client.Do(request)
		if err != nil {
			count--
			_err = err
			continue
		} else {
			httpResp = response
			break
		}
	}
	if count == 0 {
		log.Println(_err.Error())
		return _err.Error(), Agent{}
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
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
	client := http.Client{Timeout: 5 * time.Second}
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

	count := len(url)
	var _err error
	var httpResp *http.Response
	for i := 0; i < len(url); i++ {
		request, _ := http.NewRequest(http.MethodPut, url[i]+"/update",
			bytes.NewBuffer(jsonRequestMessage))
		request.Header.Add("Accept", "application/json")
		response, err := client.Do(request)
		if err != nil {
			count--
			_err = err
			continue
		} else {
			httpResp = response
			break
		}
	}
	if count == 0 {
		log.Println(_err.Error())
		return _err.Error()
	}
	defer httpResp.Body.Close()
	body, err := ioutil.ReadAll(httpResp.Body)
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

func NamesPrint(names []string) string {
	var strNames string
	name := "Nombre "
	for i, elem := range names {
		strNames += name + strconv.Itoa(i+1) + ": " + elem + "\n"
	}
	return strNames
}

func DescsPrint(descs []string) string {
	var strDescs string
	name := "Descripción "
	for i, elem := range descs {
		strDescs += name + strconv.Itoa(i+1) + ": " + elem + "\n"
	}
	return strDescs
}
