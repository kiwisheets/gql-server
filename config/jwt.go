package config

import (
	"crypto/ecdsa"
	"sync"
)

type JWTConfig struct {
	privateKey *ecdsa.PrivateKey
	mutex      sync.Mutex
}

func (j JWTConfig) GetPrivateKey() *ecdsa.PrivateKey {
	var privateKey ecdsa.PrivateKey
	j.mutex.Lock()
	privateKey = *j.privateKey
	j.mutex.Unlock()
	return &privateKey
}
