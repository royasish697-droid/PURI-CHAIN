package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Index     int    `json:"index"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
	PrevHash  string `json:"prevHash"`
	Hash      string `json:"hash"`
	Nonce     int    `json:"nonce"`
}

var Blockchain []Block

// calculate SHA256 hash for a block
func calculateHash(b Block) string {
	record := strconv.Itoa(b.Index) + b.Timestamp + b.Data + b.PrevHash + strconv.Itoa(b.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

// proof-of-work: find nonce where hash has difficulty leading zeros
func mineBlock(oldBlock Block, data string, difficulty int) Block {
	newBlock := Block{}
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Nonce = 0

	targetPrefix := strings.Repeat("0", difficulty)

	for {
		newBlock.Hash = calculateHash(newBlock)
		if strings.HasPrefix(newBlock.Hash, targetPrefix) {
			// found
			break
		}
		newBlock.Nonce++
	}
	return newBlock
}

func isChainValid(chain []Block) bool {
	for i := 1; i < len(chain); i++ {
		prev := chain[i-1]
		curr := chain[i]

		if curr.PrevHash != prev.Hash {
			return false
		}
		if calculateHash(curr) != curr.Hash {
			return false
		}
	}
	return true
}

// Add new block to global Blockchain
func addBlock(data string, difficulty int) {
	oldBlock := Blockchain[len(Blockchain)-1]
	newBlock := mineBlock(oldBlock, data, difficulty)
	Blockchain = append(Blockchain, newBlock)
	fmt.Println("Mined Block:", newBlock.Index, newBlock.Hash)
}

// HTTP handlers
type MineRequest struct {
	Data       string `json:"data"`
	Difficulty int    `json:"difficulty"`
}

func handleMine(w http.ResponseWriter, r *http.Request) {
	// only POST
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	var req MineRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	// default difficulty
	if req.Difficulty < 1 {
		req.Difficulty = 3
	}
	addBlock(req.Data, req.Difficulty)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Blockchain[len(Blockchain)-1]) // return last block
}

func handleChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Blockchain)
}

func startServer() {
	http.HandleFunc("/mine", handleMine)
	http.HandleFunc("/chain", handleChain)

	log.Println("🔥 PURI-CHAIN simple node running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createGenesisBlock() Block {
	genesis := Block{
		Index:     0,
		Timestamp: time.Now().String(),
		Data:      "Genesis Block",
		PrevHash:  "0",
		Nonce:     0,
	}
	genesis.Hash = calculateHash(genesis)
	return genesis
}

func main() {
	// init chain with genesis
	Blockchain = append(Blockchain, createGenesisBlock())
	fmt.Println("Genesis:", Blockchain[0].Hash)

	// start server
	startServer()
}
