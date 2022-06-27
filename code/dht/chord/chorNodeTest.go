package chord

import (
	"crypto/sha1"
	"fmt"
	"time"
)

func JoinTest(ip, port, ipSucc, portSucc string) {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))
	val := h.Sum(nil)
	info := NodeInfo{NodeID: val, EndPoint: Address{IP: ip, Port: port}}
	var i *NodeInfo = nil
	if ipSucc != "" {
		i = &NodeInfo{NodeID: val, EndPoint: Address{IP: ipSucc, Port: portSucc}}
	}
	n := NewNode(info, DefaultConfig(), i, getAddr(info.EndPoint))
	fmt.Println(n.getSuccessor())
	fmt.Println(n.getSuccessor())
	fmt.Println(n.Info.NodeID)
}

func SleepTest(ip, port, ipSucc, portSucc string) {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))
	val := h.Sum(nil)
	info := NodeInfo{NodeID: val, EndPoint: Address{IP: ip, Port: port}}
	var i *NodeInfo = nil
	if ipSucc != "" {
		i = &NodeInfo{NodeID: val, EndPoint: Address{IP: ipSucc, Port: portSucc}}
	}
	n := NewNode(info, DefaultConfig(), i, getAddr(info.EndPoint))
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
	time.Sleep(25 * time.Second)
	n.stabilize()
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())

}

func StabilizeTest(ip, port, ipSucc, portSucc string) {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))
	val := h.Sum(nil)
	info := NodeInfo{NodeID: val, EndPoint: Address{IP: ip, Port: port}}
	var i *NodeInfo = nil
	if ipSucc != "" {
		i = &NodeInfo{NodeID: val, EndPoint: Address{IP: ipSucc, Port: portSucc}}
	}
	n := NewNode(info, DefaultConfig(), i, getAddr(info.EndPoint))
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
	time.Sleep(15 * time.Second)
	n.stabilize()
	n.checkPredecesor()
	time.Sleep(60 * time.Second)
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
}

func FireTest(ip, port, ipSucc, portSucc string) {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))
	val := h.Sum(nil)
	info := NodeInfo{NodeID: val, EndPoint: Address{IP: ip, Port: port}}
	var n *Node
	var i *NodeInfo = nil
	if ipSucc != "" {
		i = &NodeInfo{NodeID: val, EndPoint: Address{IP: ipSucc, Port: portSucc}}
	}
	n = NewNode(info, DefaultConfig(), i, getAddr(info.EndPoint))
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
	time.Sleep(60 * time.Second)
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
}

func SelfTest(ip, port, ipSucc, portSucc string) {
	h := sha1.New()
	h.Write([]byte(ip + ":" + port))
	val := h.Sum(nil)
	info := NodeInfo{NodeID: val, EndPoint: Address{IP: ip, Port: port}}
	var n *Node
	var i *NodeInfo = nil
	if ipSucc != "" {
		i = &NodeInfo{NodeID: val, EndPoint: Address{IP: ipSucc, Port: portSucc}}
	}
	n = NewNode(info, DefaultConfig(), i, getAddr(info.EndPoint))
	time.Sleep(30 * time.Second)
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" succ", n.getSuccessor())
	fmt.Println(n.Info.EndPoint.IP+":"+n.Info.EndPoint.Port+" pred", n.getPredecessor())
}
