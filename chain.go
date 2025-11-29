package main

import (
	"errors"
	"fmt"
	"time"
)

var Blockchain []Block
var TxPool []Transaction
var Difficulty = 4 // default, changeable

func createGenesis() Block {
	gen := Block{
		Index:        0,
		Timestamp:    time.Now().Format(time.RFC3339Nano),
		Transactions: []Transaction{},
		PrevHash:     "0",
		Nonce:        0,
		Difficulty:   Difficulty,
	}
	gen.Hash = calculateHash(gen)
	return gen
}

func InitChain() {
	Blockchain = []Block{createGenesis()}
	TxPool = []Transaction{}
}

func AddBlockWithMining(minerAddr string) (Block, error) {
	// add reward transaction and mine all txs from pool
	reward := Transaction{From: "MINER", To: minerAddr, Amount: 50}
	txs := append([]Transaction{reward}, TxPool...)
	newBlock := mineBlock(Blockchain[len(Blockchain)-1], txs, Difficulty)
	Blockchain = append(Blockchain, newBlock)
	// clear pool
	TxPool = []Transaction{}
	return newBlock, nil
}

func AddTransaction(tx Transaction, pubKeyHex string) error {
	// Verify signature: need to reconstruct pub key from address? Simpler: user will send pubkey and signature
	// For now just basic checks: non-empty fields
	if tx.From == "" || tx.To == "" || tx.Amount <= 0 {
		return errors.New("invalid tx")
	}
	// Append to pool
	TxPool = append(TxPool, tx)
	return nil
}

func IsChainValid(chain []Block) bool {
	for i := 1; i < len(chain); i++ {
		prev := chain[i-1]
		cur := chain[i]
		if cur.PrevHash != prev.Hash {
			return false
		}
		if calculateHash(cur) != cur.Hash {
			return false
		}
	}
	return true
}
