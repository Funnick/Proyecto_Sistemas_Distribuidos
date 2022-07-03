package chord

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"log"
	"sync"
	"time"
)

type NodeInfo struct {
	NodeID   []byte
	EndPoint Address
}

type Node struct {
	Info NodeInfo

	cnf *Config

	ft      fingerTable
	ftMutex sync.RWMutex

	predInfo  *NodeInfo
	predMutex sync.RWMutex

	succInfo  *NodeInfo
	succMutex sync.RWMutex

	db      DataBasePlatform
	dbMutex sync.RWMutex

	stopCh chan struct{}

	next int
}

func NewNode(ip, port, dbName, knowIP, knowPort string, cnf *Config) *Node {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))

	info := NodeInfo{
		NodeID:   h.Sum(nil),
		EndPoint: Address{IP: ip, Port: port},
	}

	node := &Node{
		Info:   info,
		cnf:    cnf,
		db:     NewDataBase(dbName),
		next:   0,
		stopCh: make(chan struct{}),
	}

	node.ft = newFingerTable(&info, node.cnf.HashSize)

	var knowNode *NodeInfo = nil
	if knowIP != "" {
		h.Reset()
		h.Write([]byte(knowIP + ":" + knowPort))
		knowNode = &NodeInfo{NodeID: h.Sum(nil), EndPoint: Address{IP: knowIP, Port: knowPort}}
	}

	node.Join(knowNode)

	RunServer(NewRPCServer(node), node.Info.EndPoint, node.stopCh)

	NewBroadcastServer(node.Info.EndPoint.IP, node.stopCh)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				node.stabilize()
			case <-node.stopCh:
				ticker.Stop()
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-ticker.C:
				node.checkPredecesor()
			case <-node.stopCh:
				ticker.Stop()
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				node.fixFingers()
			case <-node.stopCh:
				ticker.Stop()
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(7 * time.Second)
		for {
			select {
			case <-ticker.C:
				node.checkSuccessor()
			case <-node.stopCh:
				ticker.Stop()
				return
			}
		}
	}()
	return node
}

// Gets a key string and returns the value
// of hash(key)
func (node *Node) getHashKey(key string) ([]byte, error) {
	// Create new Hash Func
	hashFunc := node.cnf.Hash()

	// Hash key and checks for error
	_, err := hashFunc.Write([]byte(key))
	if err != nil {
		return nil, err
	}

	// Return the value of hash(key)
	return hashFunc.Sum(nil), nil
}

// Join node n to Chord Ring
// If knowNode is nil then n is the only node
// In the Ring
func (n *Node) Join(knowNode *NodeInfo) error {
	var succ *NodeInfo = &NodeInfo{}
	if knowNode == nil {
		// N is the only node therefore
		// He's his successor
		succ.NodeID = n.Info.NodeID
		succ.EndPoint = n.Info.EndPoint
	} else {
		// There is al least one node on the Ring
		var err error
		succ, err = n.GetSuccessorOfKey(knowNode.EndPoint, n.Info.NodeID)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
	}
	n.setSuccessor(succ)

	return nil
}

// Stop
func (n *Node) Stop() {
	close(n.stopCh)
}

// Node private methods
func (n *Node) findSuccessorOfKey(key []byte) *NodeInfo {

	current := n.Info
	succ := n.getSuccessor()

	if betweenRightInlcude(current.NodeID, succ.NodeID, key) {
		return succ
	}

	cpn := n.closestPredecedingNode(key)

	if bytes.Equal(n.Info.NodeID, cpn.NodeID) {
		n.succMutex.RLock()
		defer n.succMutex.RUnlock()
		result := &NodeInfo{}
		result.NodeID = n.succInfo.NodeID
		result.EndPoint = n.succInfo.EndPoint
		return result
	}

	succOfKey, err := n.GetSuccessorOfKey(cpn.EndPoint, key)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return succOfKey
}

func (n *Node) closestPredecedingNode(key []byte) *NodeInfo {
	n.ftMutex.RLock()
	defer n.ftMutex.RUnlock()

	current := n.Info

	for i := len(n.ft.Table) - 1; i >= 0; i-- {
		ftI := n.ft.Table[i]

		if ftI == nil || ftI.SuccNode == nil {
			continue
		}
		if between(current.NodeID, key, ftI.ID) {
			var cpn *NodeInfo = &NodeInfo{}
			cpn.NodeID = ftI.SuccNode.NodeID
			cpn.EndPoint = ftI.SuccNode.EndPoint
			return cpn
		}
	}

	result := &NodeInfo{}
	result.NodeID = current.NodeID
	result.EndPoint = current.EndPoint
	return result
}

// get successor of n
func (n *Node) getSuccessor() *NodeInfo {
	var succ *NodeInfo = &NodeInfo{}

	n.succMutex.RLock()
	defer n.succMutex.RUnlock()

	if n.succInfo == nil {
		return nil
	}

	succ.NodeID = n.succInfo.NodeID
	succ.EndPoint = n.succInfo.EndPoint

	return succ
}

// get predecessor of n
func (n *Node) getPredecessor() *NodeInfo {
	var pred *NodeInfo = &NodeInfo{}

	n.predMutex.RLock()
	defer n.predMutex.RUnlock()

	if n.predInfo == nil {
		return nil
	}

	pred.NodeID = n.predInfo.NodeID
	pred.EndPoint = n.predInfo.EndPoint

	return pred
}

// set successor of n to succ
func (n *Node) setSuccessor(succ *NodeInfo) {
	n.succMutex.Lock()
	defer n.succMutex.Unlock()

	if succ == nil {
		n.succInfo = nil
		return
	}

	n.succInfo = &NodeInfo{NodeID: succ.NodeID, EndPoint: succ.EndPoint}
}

// set predecessor of n to pred
func (n *Node) setPredecessor(pred *NodeInfo) {
	n.predMutex.Lock()
	defer n.predMutex.Unlock()

	if pred == nil {
		n.predInfo = nil
		return
	}

	n.predInfo = &NodeInfo{NodeID: pred.NodeID, EndPoint: pred.EndPoint}
}

//
func (n *Node) setPosFT(pos int, node NodeInfo) {
	n.ftMutex.Lock()
	defer n.ftMutex.Unlock()

	n.ft.Table[pos].SuccNode = &node
}

// Stabilize
func (n *Node) stabilize() {
	succ := n.getSuccessor()

	if succ == nil {
		return
	}
	predOfSucc, err := n.GetPredecessorOf(succ.EndPoint)
	if err != nil {
		return
	}

	if predOfSucc == nil {
		n.Notify(succ.EndPoint)
		return
	}
	if between(n.Info.NodeID, succ.NodeID, predOfSucc.NodeID) {
		if !bytes.Equal(predOfSucc.NodeID, n.Info.NodeID) {
			n.ReplicateKey(predOfSucc.EndPoint)
		}
		n.setSuccessor(predOfSucc)
	}
	newSucc := n.getSuccessor()
	n.Notify(newSucc.EndPoint)
}

// Notify
func (n *Node) notify(newPredecessor *NodeInfo) {
	pred := n.getPredecessor()

	if pred == nil || between(pred.NodeID, n.Info.NodeID, newPredecessor.NodeID) {
		// transferir llaves
		if !bytes.Equal(newPredecessor.NodeID, n.Info.NodeID) {
			n.SendPredecessorKeys(newPredecessor.EndPoint, newPredecessor.NodeID)
		}
		n.setPredecessor(newPredecessor)
	}
}

// Check successor
func (n *Node) checkSuccessor() {
	succ := n.getSuccessor()

	if succ == nil || !n.Ping(succ.EndPoint) {
		log.Println("Buscando sucesor")
		n.ftMutex.RLock()
		defer n.ftMutex.RUnlock()
		for i := 0; i < 160; i++ {
			if n.Ping(n.ft.Table[i].SuccNode.EndPoint) {
				succ = n.ft.Table[i].SuccNode
				n.setSuccessor(succ)
				// Replicar llaves
				if !bytes.Equal(succ.NodeID, n.Info.NodeID) {
					n.ReplicateKey(succ.EndPoint)
				}
				return
			}
		}
		n.setSuccessor(&n.Info)
	}
}

// Check predecessor
func (n *Node) checkPredecesor() {
	pred := n.getPredecessor()

	if pred == nil {
		return
	}

	if !n.Ping(pred.EndPoint) {
		n.setPredecessor(nil)
	}
}

// Fix fingers
func (n *Node) fixFingers() {
	n.next = (1 + n.next) % n.cnf.HashSize

	key := calculateFingerEntryID(n.Info.NodeID, n.next, n.cnf.HashSize)
	nodeInfo := n.findSuccessorOfKey(key)
	if nodeInfo == nil {
		return
	}
	var node NodeInfo = NodeInfo{NodeID: nodeInfo.NodeID, EndPoint: nodeInfo.EndPoint}
	n.setPosFT(n.next, node)
}

// Comunication methods

// GetSuccessor -> Comunication interface implementation
func (n *Node) GetSuccessorOf(addr Address) (*NodeInfo, error) {
	nodeInfo, err := getSuccessorOf(addr)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return nodeInfo, err
}

// Ask node at address addr for the successor of key
func (n *Node) GetSuccessorOfKey(addr Address, key []byte) (*NodeInfo, error) {
	nodeInfo, err := getSuccessorOfKey(addr, key)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return nodeInfo, nil
}

// Get the predecessor of nodeInfo
func (n *Node) GetPredecessorOf(addr Address) (*NodeInfo, error) {
	nodeInfo, err := getPredecessorOf(addr)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return nodeInfo, err
}

// Notify the successor of n that n exist
func (n *Node) Notify(addr Address) error {
	var info NodeInfo = NodeInfo{}
	info = n.Info
	err := notifyNode(addr, &info)
	if err != nil {
		log.Println(err.Error())
	}
	return nil
}

// Make ping to a node
func (n *Node) Ping(addr Address) bool {
	err := ping(addr)
	return err == nil
}

// Storage Methods
func (n *Node) AskForAKey(addr Address, key []byte) ([]byte, error) {
	data, err := askForAKey(addr, key)
	if err != nil {
		log.Println(err.Error(), "AFAK")
		return nil, err
	}
	return data, nil
}

func (n *Node) SendSet(addr Address, key, data []byte) error {
	err := sendSet(addr, key, data)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func (n *Node) SendReplicate(addr Address, key []byte, data []byte) error {
	err := sendReplicate(addr, key, data)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) SendDelete(addr Address, key []byte) error {
	err := sendDelete(addr, key)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func (n *Node) SendRepDel(addr Address, key []byte) error {
	err := sendRepDel(addr, key)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func (n *Node) ReplicateKey(addr Address) error {
	fmt.Println("Obteniendo RLock")
	n.dbMutex.RLock()
	fmt.Println("RLock obtenido")

	log.Println(n.Info.EndPoint, "replicando al sucesor", addr)

	rows, err := n.db.GetKeyData()
	fmt.Println("Soltando RLock")
	n.dbMutex.RUnlock()
	fmt.Println("Suelto RLock")
	if err != nil {
		log.Println(err.Error())
		return err
	}

	fmt.Println("repl succ rows len", len(rows))
	for _, elem := range rows {
		err = n.SendReplicate(addr, elem.Key, elem.Data)
		if err != nil {
			fmt.Println("repl succ")
			log.Println(err.Error())
		}
	}

	return nil
}

func (n *Node) SendPredecessorKeys(addr Address, nID []byte) error {
<<<<<<< HEAD
	log.Println("Obteniendo RLock")
	n.dbMutex.RLock()
	log.Println("Obtenido RLock")
	rows, err := n.db.GetKeyData()
	log.Println("Soltando RLock")
	n.dbMutex.RUnlock()
	log.Println("Suelto RLock")
	pred := n.getPredecessor()
	log.Println(n.Info.EndPoint, "replicando al predecesor", addr)
=======
	fmt.Println("Obteniendo RLock")
	n.dbMutex.RLock()
	fmt.Println("RLock obtenido")

	rows, err := n.db.GetKeyData()
	fmt.Println("Soltando RLock")
	n.dbMutex.RUnlock()
	fmt.Println("Suelto RLock")

	pred := n.getPredecessor()
	log.Println(n.Info.EndPoint, "replicando al predecesor", addr)

>>>>>>> origin/DHT
	if err != nil {
		log.Println(err.Error())
		return err
	}

	fmt.Println("repl pred rows len", len(rows))
	if pred == nil {
		for _, elem := range rows {
			if bytes.Compare(elem.Key, nID) < 1 {
				n.SendReplicate(addr, elem.Key, elem.Data)
				fmt.Println("repl pred 1")
			}
		}
	} else {
		for _, elem := range rows {
			if betweenRightInlcude(pred.NodeID, nID, elem.Key) {
				err := n.SendReplicate(addr, elem.Key, elem.Data)
				if err != nil {
					log.Println(err.Error())
				}
				fmt.Println("repl pred 2")
			}
		}
	}

	return nil
}
