package chord

import (
	"crypto/sha1"
	"hash"
)

type Config struct {
	Hash     func() hash.Hash
	HashSize int
}

func DefaultConfig() *Config {
	defcng := &Config{
		Hash: sha1.New,
	}
	defcng.HashSize = defcng.Hash().Size() * 8
	return defcng
}
