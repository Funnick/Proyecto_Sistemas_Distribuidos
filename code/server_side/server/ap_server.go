package server

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

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

	// Split documentation
	docSplit := docEngine(requestMessage.Description)
	hash := sha1.New()
	hash.Write([]byte(requestMessage.IP + requestMessage.Port + requestMessage.Name))
	val := hash.Sum(nil)

	var agent Agent = Agent{Name: requestMessage.Name, AID: val, EndPoint: &endpoint,
		Password: requestMessage.Password, Description: requestMessage.Description,
		Documentation: requestMessage.Documentation, DesSplit: docSplit}
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

		hash := sha1.New()
		hash.Write([]byte(requestMessage.IP + requestMessage.Port + requestMessage.Name))
		val := hash.Sum(nil)

		if bytes.Equal(item.AID, val) {
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
	responseMessage.Message = "There is no agent with AID:" + "" // change requestMessage.AID
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
		hash := sha1.New()
		hash.Write([]byte(requestMessage.IP + requestMessage.Port + requestMessage.Name))
		val := hash.Sum(nil)
		if bytes.Equal(item.AID, val) {
			if item.Password == requestMessage.Password {
				var endpoint Address = Address{IP: requestMessage.NewIP, Port: requestMessage.NewPort}
				hash.Reset()
				hash.Write([]byte(requestMessage.NewIP + requestMessage.NewPort + requestMessage.NewName))
				newVal := hash.Sum(nil)
				var agent Agent = Agent{Name: requestMessage.Name, AID: newVal, EndPoint: &endpoint,
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

// Split documentation
func docEngine(documentation string) []string {
	docSplit := strings.Split(documentation, " ")
	return docSplit
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

// Gobinary decoder
func Decoder(file string) *[]Agent {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var agents []Agent
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&agents); err != nil {
		log.Fatal(err)
	}
	return &agents
}

// Gobinary encoder
func Encoder(agents *[]Agent) {
	f, err := os.Create("agents.gob")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	if err := enc.Encode(agents); err != nil {
		log.Fatal(err)
	}
}

// Main function
func StartServer(ip string, port string) {
	// Adding agents mock data
	agents = append(agents, Agent{AID: make([]byte, 1),
		EndPoint: &Address{IP: "127.0.0.1", Port: "8000"},
		Password: "qwer", Description: "Something1",
		Documentation: "Some documentation1"})
	agents = append(agents, Agent{AID: make([]byte, 1),
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
