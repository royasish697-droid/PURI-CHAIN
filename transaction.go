package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

type Transaction struct {
	From      string `json:"from"` // address hex of public key (or "MINER" for rewards)
	To        string `json:"to"`
	Amount    int    `json:"amount"`
	Signature string `json:"signature"` // hex r||s
}

func (tx Transaction) Hash() string {
	data := tx.From + tx.To + fmt.Sprintf("%d", tx.Amount)
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

func (tx *Transaction) Sign(priv *ecdsa.PrivateKey) error {
	if tx.From == "MINER" {
		// reward tx, no signature
		return nil
	}
	hashBytes, _ := hex.DecodeString(tx.Hash())
	r, s, err := ecdsa.Sign(nil, priv, hashBytes)
	if err != nil {
		return err
	}
	rb := r.Bytes()
	sb := s.Bytes()
	sig := append(append(make([]byte, 0), rb...), sb...)
	tx.Signature = hex.EncodeToString(sig)
	return nil
}

func (tx Transaction) Verify(pub *ecdsa.PublicKey) (bool, error) {
	if tx.From == "MINER" {
		return true, nil
	}
	if tx.Signature == "" {
		return false, errors.New("no signature")
	}
	sigBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return false, err
	}
	// signature was r||s; split into half
	half := len(sigBytes) / 2
	r := new(big.Int).SetBytes(sigBytes[:half])
	s := new(big.Int).SetBytes(sigBytes[half:])
	hashBytes, _ := hex.DecodeString(tx.Hash())
	ok := ecdsa.Verify(pub, hashBytes, r, s)
	return ok, nil
}
