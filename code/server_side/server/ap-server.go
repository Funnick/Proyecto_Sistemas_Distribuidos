package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Platform struct {
	endpoint Address
	router   mux.Router
	node     DataBasePlatform
}

func NewPlatform(ip, port string) *Platform {
	endpoint := Address{IP: ip, Port: port}
	pl := &Platform{
		endpoint: endpoint,
		router:   *mux.NewRouter(),
		node:     nil,
	}

	pl.router.HandleFunc("/ap/agents", pl.GetAgents).Methods(http.MethodGet)
	pl.router.HandleFunc("/ap/create", pl.CreateNewAgent).Methods(http.MethodPost)
	pl.router.HandleFunc("/ap/delete", pl.DeleteAgent).Methods(http.MethodDelete)
	pl.router.HandleFunc("/ap/searchbyname", pl.SearchByName).Methods(http.MethodGet)
	pl.router.HandleFunc("/ap/searchbyfunc", pl.SearchByFunc).Methods(http.MethodGet)
	pl.router.HandleFunc("/ap/update", pl.UpdateAgent).Methods(http.MethodPut)

	return pl
}

// Mocket Agents
var agents []Agent

// Get all agents
func (pl *Platform) GetAgents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}

// Create an agent
func (pl *Platform) CreateNewAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage CreateAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var endpoint Address = Address{IP: requestMessage.IP, Port: requestMessage.Port}

	var agent Agent = Agent{Name: requestMessage.Name, EndPoint: &endpoint,
		Password: requestMessage.Password, Description: requestMessage.Description,
		Documentation: requestMessage.Documentation}

	agentCode, err := json.Marshal(agent)
	if err != nil {
		log.Println(err.Error())
	}

	err = pl.node.Set(agentCode, agent.Name)
	if err != nil {
		log.Println(err.Error())
		responseMessage := ResponseMessage{Message: err.Error()}
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
	responseMessage := ResponseMessage{Message: "Agent create successfully"}
	json.NewEncoder(w).Encode(responseMessage)
}

// Delete an agent
func (pl *Platform) DeleteAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage DeleteAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var responseMessage ResponseMessage
	name, err := json.Marshal(requestMessage.Name)
	if err != nil {
		log.Println(err.Error())
		return
	}

	agent, err := pl.node.GetByName(name)
	if err != nil {
		log.Println(err.Error())
		responseMessage.Message = "There is no agent with that name:" + requestMessage.Name
		json.NewEncoder(w).Encode(responseMessage)
		return
	}

	var agentF Agent
	err = json.Unmarshal(agent, &agentF)
	if err != nil {
		log.Println(err.Error())
	}

	if requestMessage.Password != agentF.Password {
		responseMessage.Message = "Wrong password"
		json.NewEncoder(w).Encode(responseMessage)
		return
	}

	desc, err := json.Marshal(requestMessage.Description)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = pl.node.Delete(name, desc)
	if err != nil {
		log.Println(err.Error())
	} else {
		responseMessage.Message = "Agent remove successfully"
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
}

// Update an agent
// TODO
func (pl *Platform) UpdateAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage UpdateAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//TODO todo esto
	var responseMessage ResponseMessage

	name, err := json.Marshal(requestMessage.Name)
	if err != nil {
		log.Println(err.Error())
		return
	}

	agent, err := pl.node.GetByName(name)
	if err != nil {
		log.Println(err.Error())
		responseMessage.Message = "There is no agent with that name:" + requestMessage.Name
		json.NewEncoder(w).Encode(responseMessage)
		return
	}

	var agentF Agent
	err = json.Unmarshal(agent, &agentF)
	if err != nil {
		log.Println(err.Error())
	}

	if requestMessage.Password != agentF.Password {
		responseMessage.Message = "Wrong password"
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
	if requestMessage.NewDescription != agentF.Description {
		agentF.Description = requestMessage.NewDescription
	}
	if requestMessage.NewDocumentation != agentF.Documentation {
		agentF.Documentation = requestMessage.NewDocumentation
	}
	if requestMessage.NewIP != agentF.EndPoint.IP {
		agentF.EndPoint.IP = requestMessage.NewIP
	}
	if requestMessage.NewPort != agentF.EndPoint.Port {
		agentF.EndPoint.Port = requestMessage.NewPort
	}
	if requestMessage.NewPassword != agentF.Password {
		agentF.Password = requestMessage.NewPassword
	}

	newAgent, err := json.Marshal(agentF)

	if err != nil {
		log.Println(err.Error())
		return
	}
	err = pl.node.Update(newAgent)
	if err != nil {
		fmt.Println(err.Error())
	}
	responseMessage.Message = "Agent updated successfully"
	json.NewEncoder(w).Encode(responseMessage)
}

// Search an agent by functionality
func (pl *Platform) SearchByFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage SearchAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var responseMessage SearchAgentMessageResponse
	//functionality, err := json.Marshal(requestMessage.Criteria)
	agentsFound, err := pl.node.GetByFun(requestMessage.Criteria)
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(agentsFound) > 0 {
		responseMessage.Message = "These services were found"
		responseMessage.AgentFound = agentsFound
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
	responseMessage.Message = "There is no agent with that description"
	json.NewEncoder(w).Encode(responseMessage)
}

func (pl *Platform) SearchByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage SearchAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var responseMessage SearchAgentMessageResponse
	agentsFound, err := pl.node.GetByName(requestMessage.Description)
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(agentsFound) > 0 {
		responseMessage.Message = "These services were found"
		responseMessage.AgentFound = agentsFound
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
	responseMessage.Message = "There is no agent with that description"
	json.NewEncoder(w).Encode(responseMessage)
}

func (pl *Platform) Run() {
	err := http.ListenAndServe(pl.endpoint.IP+":"+pl.endpoint.Port, &pl.router)
	if err != nil {
		fmt.Println(err.Error())
	}
}
