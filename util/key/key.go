package key

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func ParseEcPrivateKeyFromPemStr(privateKeyPEM string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	p, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func ParseEcPublicKeyFromPemStr(publicKeyPEM string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	c, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch p := c.PublicKey.(type) {
	case *ecdsa.PublicKey:
		return p, nil
	default:
		break
	}
	return nil, fmt.Errorf("key type is not EC")
}
