package server

import (
	"crypto/sha1"
	"encoding/json"
	"log"
	"net/http"

	"server/chord"

	"github.com/gorilla/mux"
)

type Platform struct {
	endpoint Address
	router   mux.Router
	// Aquí creo que debería ir un chord.Node
	node chord.DBChord // ir a la línea 26
}

func NewPlatform(ip, port string) *Platform {
	endpoint := Address{IP: ip, Port: port}

	pl := &Platform{
		endpoint: endpoint,
		router:   *mux.NewRouter(),
		node:     nil, // Aquí entonces habría que inicializar el nodo y eso no sé cómo hacerlo
	}

	pl.router.HandleFunc("/ap/agents", pl.GetAgents).Methods(http.MethodGet)
	pl.router.HandleFunc("/ap/create", pl.CreateNewAgent).Methods(http.MethodPost)
	pl.router.HandleFunc("/ap/delete", pl.DeleteAgent).Methods(http.MethodDelete)
	pl.router.HandleFunc("/ap/searchbyname", pl.SearchByName).Methods(http.MethodGet)
	pl.router.HandleFunc("/ap/searchbydesc", pl.SearchByDesc).Methods(http.MethodGet)
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

	var agent Agent = Agent{Name: requestMessage.Name, EndPoint: endpoint,
		Password: requestMessage.Password, Description: requestMessage.Description,
		Documentation: requestMessage.Documentation}

	agentCode, err := json.Marshal(agent)
	if err != nil {
		log.Println(err.Error())
	}

	err = pl.node.Set(agent.Name, agent.Description, agentCode)
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
	name := requestMessage.Name
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

	desc := requestMessage.Description
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
func (pl *Platform) UpdateAgent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage UpdateAgentMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var responseMessage ResponseMessage

	name := requestMessage.Name
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
	if requestMessage.NewDescription != agentF.Description && requestMessage.NewDescription != "" {
		agentF.Description = requestMessage.NewDescription
	}
	if requestMessage.NewDocumentation != agentF.Documentation && requestMessage.NewDocumentation != "" {
		agentF.Documentation = requestMessage.NewDocumentation
	}
	if requestMessage.NewIP != agentF.EndPoint.IP && requestMessage.NewIP != "" {
		agentF.EndPoint.IP = requestMessage.NewIP
	}
	if requestMessage.NewPort != agentF.EndPoint.Port && requestMessage.NewPort != "" {
		agentF.EndPoint.Port = requestMessage.NewPort
	}
	if requestMessage.NewPassword != agentF.Password && requestMessage.NewPassword != "" {
		agentF.Password = requestMessage.NewPassword
	}

	newAgent, err := json.Marshal(agentF)

	if err != nil {
		log.Println(err.Error())
		return
	}
	err = pl.node.Update(agentF.Name, newAgent)
	if err != nil {
		log.Println(err.Error())
	}
	responseMessage.Message = "Agent updated successfully"
	json.NewEncoder(w).Encode(responseMessage)
}

// Search an agent by description
func (pl *Platform) SearchByDesc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage SearchAgentDescMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var responseMessage SearchAgentMessageResponse
	//functionality, err := json.Marshal(requestMessage.Criteria)
	agentsFound, err := pl.node.GetByName(requestMessage.Description)
	if err != nil {
		log.Println(err.Error())
		responseMessage.Message = "There is no agent with that description"
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
	var agent Agent
	json.Unmarshal(agentsFound, &agent)
	responseMessage.Message = ""
	responseMessage.AgentFound = agent
	json.NewEncoder(w).Encode(responseMessage)
}

func (pl *Platform) SearchByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestMessage SearchAgentNameMessage
	err := json.NewDecoder(r.Body).Decode(&requestMessage)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var responseMessage SearchAgentMessageResponse
	agentsFound, err := pl.node.GetByName(requestMessage.Name)
	if err != nil {
		log.Println(err.Error())
		responseMessage.Message = "There is no agent with that name"
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
	var agent Agent
	json.Unmarshal(agentsFound, &agent)
	responseMessage.Message = ""
	responseMessage.AgentFound = agent
	json.NewEncoder(w).Encode(responseMessage)
}

func (pl *Platform) Run(port string) {
	h := sha1.New()
	h.Write([]byte(pl.endpoint.IP + ":" + port))
	val := h.Sum(nil)
	info := chord.NodeInfo{NodeID: val, EndPoint: chord.Address{IP: pl.endpoint.IP, Port: port}}
	pl.node = chord.NewNode(info, chord.DefaultConfig(), nil, "pepeDB")
	err := http.ListenAndServe(pl.endpoint.IP+":"+pl.endpoint.Port, &pl.router)
	if err != nil {
		log.Println(err.Error())
	}
}
