package main

import (
	"crypto/sha1"
	"fmt"
)

func main() {
	h := sha1.New()
	s := "AgentPepe"
	h.Write([]byte(s))
	val := h.Sum(nil)
	fmt.Println(val)
	fmt.Println(len(val))
	fmt.Println(h.Size())
}
