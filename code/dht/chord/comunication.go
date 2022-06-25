package chord

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

// Request structs
type EmptyRequest struct{}
type KeyRequest struct{ Key []byte }
type NodeInfoRequest struct{ NInfo *NodeInfo }

// Response structs
type EmptyResponse struct{}
type NodeInfoResponse struct {
	NInfo *NodeInfo
	IsNil bool
}

// Alive error
type PingResquestError struct{}

func (p PingResquestError) Error() string {
	return "server off"
}

// Interface for comunication between
// Chord nodes
// All nodes have to implemet the interface
type Comunication interface {
	GetSuccessorOf(Address) (NodeInfo, error)
	GetSuccessorOfKey(Address, []byte) (NodeInfo, error)
	GetPredecessorOf(Address) (NodeInfo, error)
	Notify(Address) error
	Ping(Address) bool
}

// Server logic for rpc

// Register a service
func RegisterNodeOnRPCServer(n *Node) {
	rpc.Register(n)
}

// Run a http server
func RunRPCServer(addr Address) {
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", getAddr(addr))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	go http.Serve(listener, nil)
}

// Exported rpc method for GetSuccessor
func (n *Node) GetSuccessorRPC(request *EmptyRequest, response *NodeInfoResponse) error {
	succ := n.getSuccessor()
	if succ == nil {
		response.NInfo = &NodeInfo{}
		response.IsNil = true
		return nil
	}

	response.NInfo = &NodeInfo{}
	response.NInfo.NodeID = succ.NodeID
	response.NInfo.EndPoint = succ.EndPoint
	response.IsNil = false
	return nil
}

// Exported rpc method for GetSuccessorOfKey
func (n *Node) GetSuccessorOfKeyRPC(resquest *KeyRequest, response *NodeInfoResponse) error {
	info := n.findSuccessorOfKey(resquest.Key)
	if info == nil {
		response.NInfo = &NodeInfo{}
		response.IsNil = true
		return nil
	}

	response.NInfo = &NodeInfo{}
	response.NInfo.NodeID = info.NodeID
	response.NInfo.EndPoint = info.EndPoint
	response.IsNil = false
	return nil
}

// Exported rpc method for getPredecessor
func (n *Node) GetPredecessorOfRPC(request *EmptyRequest, response *NodeInfoResponse) error {
	predInfo := n.getPredecessor()
	if predInfo == nil {
		response.NInfo = &NodeInfo{}
		response.IsNil = true
		return nil
	}

	response.NInfo = &NodeInfo{}
	response.NInfo.NodeID = predInfo.NodeID
	response.NInfo.EndPoint = predInfo.EndPoint
	response.IsNil = false
	return nil
}

// Exported rpc method for notify
func (n *Node) NotifyNode(request *NodeInfoRequest, response *EmptyResponse) error {
	n.notify(request.NInfo)
	return nil
}

func (n *Node) MakePingRPC(request *EmptyRequest, response *EmptyResponse) error {
	return nil
}

// Client logic for rpc

// Connect the client with server at address addr
// And calls GetSuccessor of node at address addr
func getSuccessorOf(addr Address) (*NodeInfo, error) {
	client, err := rpc.DialHTTP("tcp", getAddr(addr))
	if err != nil {
		return nil, err
	}

	var response *NodeInfoResponse = &NodeInfoResponse{}
	err = client.Call("Node.GetSuccessorRPC", &EmptyRequest{}, response)
	if err != nil {
		return nil, err
	}

	if response.IsNil {
		return nil, nil
	}

	return response.NInfo, nil
}

func getSuccessorOfKey(addr Address, key []byte) (*NodeInfo, error) {
	client, err := rpc.DialHTTP("tcp", getAddr(addr))
	if err != nil {
		return nil, err
	}
	var response *NodeInfoResponse = &NodeInfoResponse{}
	err = client.Call("Node.GetSuccessorOfKeyRPC", &KeyRequest{Key: key}, response)
	if err != nil {
		return nil, err
	}

	if response.IsNil {
		return nil, nil
	}

	return response.NInfo, nil
}

// Ask the node at address addr for his predecessor
func getPredecessorOf(addr Address) (*NodeInfo, error) {
	client, err := rpc.DialHTTP("tcp", getAddr(addr))
	if err != nil {
		return nil, err
	}

	var response *NodeInfoResponse = &NodeInfoResponse{}
	err = client.Call("Node.GetPredecessorOfRPC", &EmptyRequest{}, response)
	if err != nil {
		return nil, err
	}

	if response.IsNil {
		return nil, nil
	}

	return response.NInfo, nil
}

func notifyNode(addr Address, n *NodeInfo) error {
	client, err := rpc.DialHTTP("tcp", getAddr(addr))
	if err != nil {
		return err
	}

	err = client.Call("Node.NotifyNode", &NodeInfoRequest{NInfo: n}, nil)
	return err
}

func ping(addr Address) error {
	_, err := rpc.DialHTTP("tcp", getAddr(addr))
	if err != nil {
		return PingResquestError{}
	}

	return nil
}
