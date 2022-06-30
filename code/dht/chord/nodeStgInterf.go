package chord

type DBChord interface {
	GetByName(string) (string, error)
	GetByFun(string) ([]string, error)
	Set(string, string) error
	Update(string, string) error
	Delete(string) error
}

const (
	Name string = "Name"
	Fun  string = "Function"
)

func (n *Node) GetByName(agentName string) (string, error) {
	key, err := n.getHashKey(agentName)
	if err != nil {
		return "", err
	}

	nInfo := n.findSuccessorOfKey(key)
	return n.AskForAKey(nInfo.EndPoint, key)
}

/*
func (n *Node) GetByFun(fun string) ([]string, error) {
	key, err := n.getHashKey(Fun+fun)
}
*/

func (n *Node) Set(name string, data string) error {
	key, err := n.getHashKey(name)
	if err != nil {
		return err
	}

	nInfo := n.findSuccessorOfKey(key)
	return n.SendSet(nInfo.EndPoint, key, data)
}

func (n *Node) Update(name string, data string) error {
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

func (n *Node) Delete(name string) error {
	key, err := n.getHashKey(name)
	if err != nil {
		return err
	}

	nInfo := n.findSuccessorOfKey(key)
	return n.SendDelete(nInfo.EndPoint, key)
}
