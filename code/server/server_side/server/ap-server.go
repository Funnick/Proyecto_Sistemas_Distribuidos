package server

import (
	"crypto/sha1"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"server/chord"

	"github.com/gorilla/mux"
)

type Platform struct {
	endpoint Address
	router   mux.Router
	node     chord.DBChord
}

func NewPlatform(ip, port string) *Platform {
	endpoint := Address{IP: ip, Port: port}

	pl := &Platform{
		endpoint: endpoint,
		router:   *mux.NewRouter(),
		node:     nil,
	}

	pl.router.HandleFunc("/ap/names", pl.GetAgentsNames).Methods(http.MethodGet)
	pl.router.HandleFunc("/ap/descriptions", pl.GetAgentsDescs).Methods(http.MethodGet)
	pl.router.HandleFunc("/ap/create", pl.CreateNewAgent).Methods(http.MethodPost)
	pl.router.HandleFunc("/ap/delete", pl.DeleteAgent).Methods(http.MethodDelete)
	pl.router.HandleFunc("/ap/searchbyname", pl.SearchByName).Methods(http.MethodGet)
	pl.router.HandleFunc("/ap/searchbydesc", pl.SearchByDesc).Methods(http.MethodGet)
	pl.router.HandleFunc("/ap/update", pl.UpdateAgent).Methods(http.MethodPut)

	return pl
}

// Stop
func (pl *Platform) Stop() {
	pl.node.Stop()
}

// Get all agents
func (pl *Platform) GetAgentsNames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var responseMessage GetAllResponse
	responseFound, err := pl.node.GetAllNames()
	if err != nil {
		log.Println(err.Error())
		responseMessage.Message = "No hay agentes registrados"
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
	var resp []string
	json.Unmarshal(responseFound, &resp)
	responseMessage.Message = ""
	responseMessage.ResponsesFound = resp
	json.NewEncoder(w).Encode(responseMessage)
}

func (pl *Platform) GetAgentsDescs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var responseMessage GetAllResponse
	responseFound, err := pl.node.GetAllFun()
	if err != nil {
		log.Println(err.Error())
		responseMessage.Message = "No hay agentes registrados"
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
	var resp []string
	json.Unmarshal(responseFound, &resp)
	responseMessage.Message = ""
	responseMessage.ResponsesFound = resp
	json.NewEncoder(w).Encode(responseMessage)
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
	_, err = pl.node.GetByName(requestMessage.Name)
	if err != nil && err.Error() != "Recurso no encontrado" {
		log.Println(err.Error())
		return
	}
	if err == nil {
		responseMessage := ResponseMessage{Message: "Ya existe el agente"}
		json.NewEncoder(w).Encode(responseMessage)
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
	responseMessage := ResponseMessage{Message: "Agente creado satisfactoriamente"}
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
		responseMessage.Message = "No existe el agente:" + requestMessage.Name
		json.NewEncoder(w).Encode(responseMessage)
		return
	}

	var agentF Agent
	err = json.Unmarshal(agent, &agentF)
	if err != nil {
		log.Println(err.Error())
		return
	}

	if requestMessage.Password != agentF.Password {
		responseMessage.Message = "Contraseña incorrecta"
		json.NewEncoder(w).Encode(responseMessage)
		return
	}

	err = pl.node.Delete(name, agentF.Description)
	if err != nil {
		log.Println(err.Error())
	} else {
		responseMessage.Message = "Agente removido satisfactoriamente"
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
		responseMessage.Message = "No existe el agente:" + requestMessage.Name
		json.NewEncoder(w).Encode(responseMessage)
		return
	}

	var agentF Agent
	err = json.Unmarshal(agent, &agentF)
	if err != nil {
		log.Println(err.Error())
	}
	desc := agentF.Description

	if requestMessage.Password != agentF.Password {
		responseMessage.Message = "Contraseña incorrecta"
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
	err = pl.node.Update(agentF.Name, desc, agentF.Description, newAgent)
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
	agentsFound, err := pl.node.GetByFun(requestMessage.Description)
	if err != nil {
		log.Println(err.Error())
		responseMessage.Message = "No existen agentes con esa descripción"
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
		responseMessage.Message = "No existe el agente"
		json.NewEncoder(w).Encode(responseMessage)
		return
	}
	var agent Agent
	json.Unmarshal(agentsFound, &agent)
	responseMessage.Message = ""
	responseMessage.AgentFound = agent
	json.NewEncoder(w).Encode(responseMessage)
}

func (pl *Platform) Join(ip, port string) error {
	if ip == "" {
		ip = pl.Multicast(ip)
	}
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))

	knowNode := &chord.NodeInfo{
		NodeID:   h.Sum(nil),
		EndPoint: chord.Address{IP: ip, Port: port},
	}

	return pl.node.Join(knowNode)
}

func (pl *Platform) Multicast(knowIP string) string {
	if knowIP == "" {
		ip := pl.endpoint.IP
		ipSplit := strings.Split(ip, ".")[:3]
		var subnet string
		for _, value := range ipSplit {
			subnet += value + "."
		}
		for i := 1; i < 10; i++ {
			tryIP := subnet + strconv.Itoa(i)
			client, err := net.Dial("tcp", tryIP+":6002")
			if err != nil {
				continue
			}
			client.Close()

			knowIP = tryIP

			break
		}
	}
	return knowIP
}

func (pl *Platform) Run(port, nameDB, knowIP string) {
	knowIP = pl.Multicast(knowIP)
	pl.node = chord.NewNode(pl.endpoint.IP, port, nameDB, knowIP, "6001", chord.DefaultConfig())
	err := http.ListenAndServe(pl.endpoint.IP+":"+pl.endpoint.Port, &pl.router)
	if err != nil {
		log.Println(err.Error())
	}
}
