package chord

import (
	"sync"
)

type Address struct {
	IP   string
	Port string
}

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

	// TODO
}

// func newNode

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
func (n *Node) join(knowNode *Node) error {
	if knowNode == nil {
		// N is the only node therefore
		// He's his successor
		n.succMutex.Lock()
		n.succInfo = &n.Info
		n.succMutex.Unlock()

		return nil
	}

	// There is al least one node on the Ring
	
}

// func (n *Node) Find(key string) (*Node, error) {}

// func (n *Node) Get(key string) ([]bytes, error) {}

// func (n *Node) Set(key, value string) error {}

// func (n *Node) Delete(key string) error {}
