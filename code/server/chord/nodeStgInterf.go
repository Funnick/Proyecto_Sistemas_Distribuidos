package chord

import (
	"encoding/json"
	"log"
)

type DBChord interface {
	GetByName(string) ([]byte, error)
	GetByFun(string) ([][]byte, error)
	Set(string, string, []byte) error
	Update(string, string, string, []byte) error
	Delete(string, string) error
	GetAllNames() ([]byte, error)
	GetAllFun() ([]byte, error)
	Stop()
	Join(*NodeInfo) error
}

const (
	Names string = "Names"
	Funs  string = "Functions"
)

func (n *Node) GetAllNames() ([]byte, error) {
	key, err := n.getHashKey(Names)
	if err != nil {
		return nil, err
	}

	nInfo := n.findSuccessorOfKey(key)
	return n.AskForAKey(nInfo.EndPoint, key)
}

func (n *Node) GetAllFun() ([]byte, error) {
	key, err := n.getHashKey(Funs)
	if err != nil {
		return nil, err
	}

	nInfo := n.findSuccessorOfKey(key)
	agentsFun, err := n.AskForAKey(nInfo.EndPoint, key)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	var af map[string][]string

	err = json.Unmarshal(agentsFun, &af)
	if err != nil {
		return nil, err
	}

	var funs []string
	for k := range af {
		funs = append(funs, k)
	}

	return json.Marshal(funs)
}

func (n *Node) GetByName(agentName string) ([]byte, error) {
	key, err := n.getHashKey(agentName)
	if err != nil {
		return nil, err
	}

	nInfo := n.findSuccessorOfKey(key)
	return n.AskForAKey(nInfo.EndPoint, key)
}

// arreglar si no existe el map de funciones
func (n *Node) GetByFun(fun string) ([][]byte, error) {
	key, err := n.getHashKey(Funs)
	if err != nil {
		return make([][]byte, 0), err
	}
	nInfo := n.findSuccessorOfKey(key)
	agentsFun, err := n.AskForAKey(nInfo.EndPoint, key)
	if err != nil {
		return make([][]byte, 0), err
	}

	// Arreglar agentsFun len 0

	var af map[string][]string

	err = json.Unmarshal(agentsFun, &af)
	if err != nil {
		return nil, err
	}

	agents := make([][]byte, 0)
	agentNames := af[fun]
	for i := range agentNames {
		key, err = n.getHashKey(agentNames[i])
		if err != nil {
			return make([][]byte, 0), err
		}
		nInfo = n.findSuccessorOfKey(key)
		agent, err := n.AskForAKey(nInfo.EndPoint, key)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		agents = append(agents, agent)
	}
	return agents, nil
}

func setNames(agentNames []byte, name string) ([]byte, error) {
	var an []string

	if len(agentNames) == 0 {
		an = make([]string, 0)
	} else {
		err := json.Unmarshal(agentNames, &an)
		if err != nil {
			return nil, err
		}
	}

	an = append(an, name)
	return json.Marshal(an)
}

func setFunctions(agentsFun []byte, fun string, name string) ([]byte, error) {
	var af map[string][]string

	if len(agentsFun) == 0 {
		af = make(map[string][]string, 1)
	} else {
		err := json.Unmarshal(agentsFun, &af)
		if err != nil {
			return nil, err
		}
	}

	// kmp
	var check bool = false
	for k := range af {
		if k == fun {
			agentNames := af[fun]
			agentNames = append(agentNames, name)
			af[fun] = agentNames
			check = true
			break
		}
	}
	if !check {
		agentNames := make([]string, 1)
		agentNames[0] = name
		af[fun] = agentNames
	}

	return json.Marshal(af)
}

func (n *Node) Set(name string, fun string, data []byte) error {
	// Actualiza la lista de los nombres de los agentes
	key, err := n.getHashKey(Names)
	if err != nil {
		return err
	}
	nInfo := n.findSuccessorOfKey(key)
	var agentNames []byte
	agentNames, err = n.AskForAKey(nInfo.EndPoint, key)
	if err != nil && err.Error() != "Recurso no encontrado" {
		return err
	}
	if err != nil {
		agentNames = make([]byte, 0)
	} else {
		n.SendDelete(nInfo.EndPoint, key)
	}
	d, err := setNames(agentNames, name)
	if err != nil {
		return err
	}
	err = n.SendSet(nInfo.EndPoint, key, d)
	if err != nil {
		return err
	}

	// Actualiza el diccionario funcion-agentes
	key, err = n.getHashKey(Funs)
	if err != nil {
		return err
	}
	nInfo = n.findSuccessorOfKey(key)
	functionsAgents, err := n.AskForAKey(nInfo.EndPoint, key)
	if err != nil && err.Error() != "Recurso no encontrado" {
		return err
	}
	if err != nil {
		functionsAgents = make([]byte, 0)
	} else {
		n.SendDelete(nInfo.EndPoint, key)
	}
	d, err = setFunctions(functionsAgents, fun, name)
	if err != nil {
		return err
	}

	err = n.SendSet(nInfo.EndPoint, key, d)
	if err != nil {
		return err
	}

	// Guarda el agente en el DHT
	key, err = n.getHashKey(name)
	nInfo = n.findSuccessorOfKey(key)
	if err != nil {
		return err
	}

	return n.SendSet(nInfo.EndPoint, key, data)
}

// Arreglar
func (n *Node) Update(name, oldFun, newFun string, data []byte) error {
	if oldFun != newFun {
		key, err := n.getHashKey(Funs)
		if err != nil {
			return err
		}

		nInfo := n.findSuccessorOfKey(key)
		agentsFun, err := n.AskForAKey(nInfo.EndPoint, key)
		if err != nil {
			return err
		}
		agentsFun, err = setFunctions(agentsFun, newFun, name)
		if err != nil {
			return err
		}
		agentsFun, err = deleteFun(agentsFun, oldFun, name)
		if err != nil {
			return err
		}
		err = n.SendDelete(nInfo.EndPoint, key)
		if err != nil {
			return err
		}
		err = n.SendSet(nInfo.EndPoint, key, agentsFun)
		if err != nil {
			return err
		}
	}

	key, err := n.getHashKey(name)
	if err != nil {
		return err
	}

	nInfo := n.findSuccessorOfKey(key)
	err = n.SendDelete(nInfo.EndPoint, key)
	if err != nil {
		return err
	}
	return n.SendSet(nInfo.EndPoint, key, data)
}

func deleteAgentName(agentNames []byte, name string) ([]byte, error) {
	var an []string
	err := json.Unmarshal(agentNames, &an)
	if err != nil {
		return nil, err
	}

	for i, elem := range an {
		if elem == name {
			an = append(an[:i], an[i+1:]...)
			break
		}
	}

	return json.Marshal(an)
}

func deleteFun(agentsFun []byte, fun, name string) ([]byte, error) {
	var af map[string][]string
	err := json.Unmarshal(agentsFun, &af)
	if err != nil {
		return nil, err
	}

	// kmp
	for k := range af {
		if k == fun {
			agentNames := af[fun]
			if len(agentNames) == 1 {
				delete(af, fun)
			} else {
				for i := range agentNames {
					if agentNames[i] == name {
						agentNames = append(agentNames[:i], agentNames[i+1:]...)
						af[fun] = agentNames
						break
					}
				}
			}
			break
		}
	}

	return json.Marshal(af)
}

func (n *Node) Delete(name string, fun string) error {
	key, err := n.getHashKey(Names)
	if err != nil {
		return err
	}
	nInfo := n.findSuccessorOfKey(key)
	agentNames, err := n.AskForAKey(nInfo.EndPoint, key)
	if err != nil && err.Error() != "Recurso no encontrado" {
		return err
	}
	if err == nil {
		agentNames, err = deleteAgentName(agentNames, name)
		if err != nil {
			return err
		}
		err = n.SendDelete(nInfo.EndPoint, key)
		if err != nil {
			return err
		}
		err = n.SendSet(nInfo.EndPoint, key, agentNames)
		if err != nil {
			return err
		}
	}

	key, err = n.getHashKey(Funs)
	if err != nil {
		return err
	}
	agentsFun, err := n.AskForAKey(nInfo.EndPoint, key)
	if err != nil && err.Error() != "Recurso no encontrado" {
		return err
	}
	if err == nil {
		agentsFun, err = deleteFun(agentsFun, fun, name)
		if err != nil {
			return err
		}
		err = n.SendDelete(nInfo.EndPoint, key)
		if err != nil {
			return err
		}
		err = n.SendSet(nInfo.EndPoint, key, agentsFun)
		if err != nil {
			return err
		}
	}

	key, err = n.getHashKey(name)
	if err != nil {
		return err
	}

	nInfo = n.findSuccessorOfKey(key)
	return n.SendDelete(nInfo.EndPoint, key)
}
