package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Platform struct {
	endpoint Address
	router   mux.Router
	node     DBChord
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
	pl.router.HandleFunc("/ap/search", pl.SearchAgent).Methods(http.MethodGet)
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
		fmt.Println(err.Error())
		return
	}
	var endpoint Address = Address{IP: requestMessage.IP, Port: requestMessage.Port}

	var agent Agent = Agent{Name: requestMessage.Name, EndPoint: &endpoint,
		Password: requestMessage.Password, Description: requestMessage.Description,
		Documentation: requestMessage.Documentation}

	err = pl.node.Set(agent)
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
		return
	}
	var responseMessage ResponseMessage
	err = pl.node.Delete("", "")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		responseMessage.Message = "Agent remove successfully"
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
	// responseMessage.Message = "Wrong password"

	responseMessage.Message = "There is no agent with that name:" + requestMessage.Name
	json.NewEncoder(w).Encode(responseMessage)
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
	var responseMessage ResponseMessage
	agent, err := pl.node.GetByName(requestMessage.Name)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = pl.node.Update("", "", string(agent))
	if err != nil {
		fmt.Println(err.Error())
	}
	responseMessage.Message = "Agent updated successfully"
	//responseMessage.Message = "Wrong password"
	json.NewEncoder(w).Encode(responseMessage)
}

// Search an agent
func (pl *Platform) SearchAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage SearchAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var responseMessage SearchAgentMessageResponse
	agentsFound, err := pl.node.GetByFunction(requestMessage.Description)
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(agentsFound) > 0 {
		responseMessage.Message = ""
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
