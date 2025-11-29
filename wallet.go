package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
	Address    string
}

func NewWallet() (*Wallet, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		return nil, err
	}
	pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	// simple address: sha256(pub) hex
	addrHash := sha256.Sum256(pub)
	addr := hex.EncodeToString(addrHash[:])
	return &Wallet{PrivateKey: priv, PublicKey: pub, Address: addr}, nil
}

func WalletFromPubkeyBytes(pub []byte) (string, error) {
	if len(pub) == 0 {
		return "", errors.New("empty pub")
	}
	addr := sha256.Sum256(pub)
	return hex.EncodeToString(addr[:]), nil
}

func PublicKeyToECDSAPub(pub []byte) (*ecdsa.PublicKey, error) {
	if len(pub)%2 != 0 {
		return nil, errors.New("invalid pubkey length")
	}
	half := len(pub) / 2
	x := new(big.Int).SetBytes(pub[:half])
	y := new(big.Int).SetBytes(pub[half:])
	return &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}, nil
}
