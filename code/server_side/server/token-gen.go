package server

import "time"

type Token struct {
	value      string
	expiration string
}

func (token *Token) Generate(name string) {
	token.value = name
	token.expiration = time.Now().String()
}
