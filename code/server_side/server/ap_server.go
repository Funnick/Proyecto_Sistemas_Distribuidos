package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ServerAP struct {
	endpoint *Address
	router   mux.Router
}

// Mocket Agents
var agents []Agent

// Handlers Functions

// Get all agents
func getAgents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}

// Create an agent
func createNewAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage CreateAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		log.Fatal(err)
		return
	}
	var endpoint Address = Address{IP: requestMessage.IP, Port: requestMessage.Port}
	var agent Agent = Agent{AID: "4", EndPoint: &endpoint,
		Password: requestMessage.Password, Description: requestMessage.Description,
		Documentation: requestMessage.Documentation}
	agents = append(agents, agent)
	responseMessage := ResponseMessage{Message: "Agent create successfully"}
	json.NewEncoder(w).Encode(responseMessage)
}

// Delete an agent
func deleteAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage DeleteAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		log.Fatal(err)
		return
	}
	var responseMessage ResponseMessage
	for i, item := range agents {
		if item.AID == requestMessage.AID {
			if item.Password == requestMessage.Password {
				agents = append(agents[:i], agents[i+1:]...)
				responseMessage.Message = "Agent remove successfully"
			} else {
				responseMessage.Message = "Wrong password"
			}
			json.NewEncoder(w).Encode(responseMessage)
			return
		}
	}
	responseMessage.Message = "There is no agent with AID:" + requestMessage.AID
	json.NewEncoder(w).Encode(responseMessage)
}

// Update an agent
// TODO
func updateAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage UpdateAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		log.Fatal(err)
		return
	}
	var responseMessage ResponseMessage
	for i, item := range agents {
		if item.AID == requestMessage.AID {
			if item.Password == requestMessage.Password {
				var endpoint Address = Address{IP: requestMessage.NewIP, Port: requestMessage.NewPort}
				var agent Agent = Agent{AID: requestMessage.AID, EndPoint: &endpoint,
					Password: requestMessage.NewPassword, Description: requestMessage.NewDescription,
					Documentation: requestMessage.NewDocumentation}
				agents[i] = agent
				responseMessage.Message = "Agent updated successfully"
			} else {
				responseMessage.Message = "Wrong password"
			}
			json.NewEncoder(w).Encode(responseMessage)
			return
		}
	}
	responseMessage.Message = "There is no agent with that ID"
	json.NewEncoder(w).Encode(responseMessage)
}

// Search an agent
func searchAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage SearchAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		log.Fatal(err)
		return
	}
	var responseMessage SearchAgentMessageResponse
	for _, item := range agents {
		if item.Description == requestMessage.Description {
			responseMessage.Message = ""
			responseMessage.AgentFound = item
			json.NewEncoder(w).Encode(responseMessage)
			return
		}
	}
	responseMessage.Message = "There is no agent with that description"
	json.NewEncoder(w).Encode(responseMessage)
}

func NewServerAP(endpoint *Address) *ServerAP {
	newServer := &ServerAP{endpoint: endpoint, router: *mux.NewRouter()}

	newServer.router.HandleFunc("/ap/agents", getAgents).Methods(http.MethodGet)
	newServer.router.HandleFunc("/ap/create", createNewAgent).Methods(http.MethodPost)
	newServer.router.HandleFunc("/ap/delete", deleteAgent).Methods(http.MethodDelete)
	newServer.router.HandleFunc("/ap/search", searchAgent).Methods(http.MethodGet)
	newServer.router.HandleFunc("/ap/update", updateAgent).Methods(http.MethodPut)

	return newServer
}

func (s *ServerAP) Run() {
	log.Fatal(http.ListenAndServe(s.endpoint.IP+":"+s.endpoint.Port, &s.router))
}

// Main function
func StartServer(ip string, port string) {
	// Adding agents mock data
	agents = append(agents, Agent{AID: "1",
		EndPoint: &Address{IP: "127.0.0.1", Port: "8000"},
		Password: "qwer", Description: "Something1",
		Documentation: "Some documentation1"})
	agents = append(agents, Agent{AID: "2",
		EndPoint: &Address{IP: "192.168.0.1", Port: "5000"},
		Password: "rewq", Description: "Something2",
		Documentation: "Some documentation2"})

	// Init Router
	router := mux.NewRouter()

	// Route Handlers / EndPoints
	router.HandleFunc("/ap/agents", getAgents).Methods(http.MethodGet)
	router.HandleFunc("/ap/create", createNewAgent).Methods(http.MethodPost)
	router.HandleFunc("/ap/delete", deleteAgent).Methods(http.MethodDelete)
	router.HandleFunc("/ap/search", searchAgent).Methods(http.MethodGet)
	router.HandleFunc("/ap/update", updateAgent).Methods(http.MethodPut)

	// Init Server
	err := http.ListenAndServe(ip+":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
