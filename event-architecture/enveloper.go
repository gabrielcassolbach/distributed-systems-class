package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
)

type Signer struct {
	PrivateKey *rsa.PrivateKey
}

func NewSigner() (*Signer, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &Signer{PrivateKey: priv}, nil
}

func (s *Signer) Sign(data interface{}) (string) {
	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)

	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, s.PrivateKey, crypto.SHA256, hash[:])
	if err != nil {
		return ""
	}

	pubBytes, _ := x509.MarshalPKIXPublicKey(&s.PrivateKey.PublicKey)
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: pubBytes})

	envelope := map[string]interface{}{
		"signature":  base64.StdEncoding.EncodeToString(sigBytes),
		"public_key": string(pubPEM),
		"data":       data,
	}

	envelopeJSON, _ := json.Marshal(envelope)
	return string(envelopeJSON)
}

func (s *Signer) Open(envelopeJSON string) (string, error) {
	var envelope map[string]interface{}
	if err := json.Unmarshal([]byte(envelopeJSON), &envelope); err != nil {
		return "", fmt.Errorf("erro ao ler o JSON: %v", err)
	}

	data := envelope["data"]
	signatureBase64 := envelope["signature"].(string)
	publicKeyPEM := envelope["public_key"].(string)

	jsonData, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonData)

	sigBytes, _ := base64.StdEncoding.DecodeString(signatureBase64)

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", fmt.Errorf("chave pública inválida")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	pubKey := pubInterface.(*rsa.PublicKey)

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], sigBytes)
	if err != nil {
		return "", fmt.Errorf("assinatura INVÁLIDA: o dado foi alterado ou a chave não bate")
	}

	return data.(string), nil
}

