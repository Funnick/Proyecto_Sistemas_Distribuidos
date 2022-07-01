package chord

import (
	"encoding/json"
)

type DBChord interface {
	GetByName(string) ([]byte, error)
	//GetByFun(string) ([]string, error)
	Set(string, []byte) error
	Update(string, []byte) error
	Delete(string) error
}

const (
	Names string = "Names"
	Funs  string = "Functions"
)

func (n *Node) GetByName(agentName string) ([]byte, error) {
	key, err := n.getHashKey(agentName)
	if err != nil {
		return nil, err
	}

	nInfo := n.findSuccessorOfKey(key)
	return n.AskForAKey(nInfo.EndPoint, key)
}

/*
func (n *Node) GetByFun(fun string) ([]string, error) {
	key, err := n.getHashKey(Fun+fun)
}
*/

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

// Hacer
func setFunctions(agentsFun []byte, fun string) ([]byte, error) {
	var af map[string][]string
	err := json.Unmarshal(agentsFun, &af)
	if err != nil {
		return nil, err
	}

	// kmp

	// Cambiar
	return nil, nil
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
	if err != nil && err.Error() != "There is no agent with that name" {
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
	/*
		// Actualiza el diccionario funcion-agentes
		key, err = n.getHashKey(Funs)
		if err != nil {
			return err
		}
		nInfo = n.findSuccessorOfKey(key)
		functionsAgents, err := n.AskForAKey(nInfo.EndPoint, key)
		if err != nil && err.Error() != "There is no agent with that name" {
			return err
		}
		if err != nil {
			functionsAgents = make([]byte, 0)
		} else {
			n.SendDelete(nInfo.EndPoint, key)
		}
		d, err = setFunctions(functionsAgents, fun)
		if err != nil {
			return err
		}

		err = n.SendSet(nInfo.EndPoint, key, d)
	*/
	// Guarda el agente en el DHT
	key, err = n.getHashKey(name)
	nInfo = n.findSuccessorOfKey(key)
	if err != nil {
		return err
	}

	return n.SendSet(nInfo.EndPoint, key, data)
}

func (n *Node) Update(name string, data []byte) error {
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

func (n *Node) Delete(name string) error {
	key, err := n.getHashKey(Names)
	nInfo := n.findSuccessorOfKey(key)
	agentNames, err := n.AskForAKey(nInfo.EndPoint, key)
	if err != nil {
		return err
	}
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

	// Falta arreglar las funciones

	key, err = n.getHashKey(name)
	if err != nil {
		return err
	}

	nInfo = n.findSuccessorOfKey(key)
	return n.SendDelete(nInfo.EndPoint, key)
}
