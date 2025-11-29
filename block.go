package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Block struct {
	Index        int            `json:"index"`
	Timestamp    string         `json:"timestamp"`
	Transactions []Transaction  `json:"transactions"`
	PrevHash     string         `json:"prevHash"`
	Hash         string         `json:"hash"`
	Nonce        int            `json:"nonce"`
	Difficulty   int            `json:"difficulty"`
}

func calculateHash(b Block) string {
	record := fmt.Sprintf("%d%s%s%s%d", b.Index, b.Timestamp, serializeTxs(b.Transactions), b.PrevHash, b.Nonce)
	h := sha256.Sum256([]byte(record))
	return hex.EncodeToString(h[:])
}

func serializeTxs(txs []Transaction) string {
	// simple serialization for hashing
	res := ""
	for _, t := range txs {
		res += t.From + t.To + fmt.Sprintf("%d", t.Amount) + t.Signature
	}
	return res
}

func mineBlock(old Block, txs []Transaction, difficulty int) Block {
	newBlock := Block{}
	newBlock.Index = old.Index + 1
	newBlock.Timestamp = time.Now().Format(time.RFC3339Nano)
	newBlock.Transactions = txs
	newBlock.PrevHash = old.Hash
	newBlock.Difficulty = difficulty
	newBlock.Nonce = 0

	prefix := strings.Repeat("0", difficulty)
	for {
		newBlock.Hash = calculateHash(newBlock)
		if strings.HasPrefix(newBlock.Hash, prefix) {
			break
		}
		newBlock.Nonce++
	}
	return newBlock
}
