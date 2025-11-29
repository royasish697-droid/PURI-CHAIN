package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func handleMine(w http.ResponseWriter, r *http.Request) {
	// miner address from URL param or body
	miner := r.URL.Query().Get("miner")
	if miner == "" {
		http.Error(w, "miner address required, use ?miner=ADDRESS", http.StatusBadRequest)
		return
	}
	block, err := AddBlockWithMining(miner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := map[string]interface{}{
		"message": "New block mined",
		"block":   block,
	}
	json.NewEncoder(w).Encode(resp)
}

func handleChain(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Blockchain)
}

func handleNewWallet(w http.ResponseWriter, r *http.Request) {
	wallet, err := NewWallet()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// return address and pubkey hex
	resp := map[string]string{
		"address": wallet.Address,
		"pubkey":  hex.EncodeToString(wallet.PublicKey),
	}
	json.NewEncoder(w).Encode(resp)
}

func handleNewTx(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var tx Transaction
	if err := json.Unmarshal(body, &tx); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	// basic add (signature verification left for client)
	if tx.From != "MINER" && tx.Signature != "" {
		// try verify: client must provide pubkey in header "X-Pubkey"
		pubHex := r.Header.Get("X-Pubkey")
		if pubHex == "" {
			http.Error(w, "missing X-Pubkey header (hex)", http.StatusBadRequest)
			return
		}
		pubBytes, err := hex.DecodeString(pubHex)
		if err != nil {
			http.Error(w, "invalid pubkey hex", http.StatusBadRequest)
			return
		}
		// reconstruct ecdsa pub
		half := len(pubBytes) / 2
		x := new(big.Int).SetBytes(pubBytes[:half])
		y := new(big.Int).SetBytes(pubBytes[half:])
		pub := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

		ok, err := tx.Verify(pub)
		if err != nil || !ok {
			http.Error(w, "signature verify failed", http.StatusBadRequest)
			return
		}
	}
	if err := AddTransaction(tx, ""); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "tx added to pool"})
}

func handleTxPool(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(TxPool)
}

func startServer() {
	http.HandleFunc("/mine", handleMine)        // GET ?miner=ADDR
	http.HandleFunc("/chain", handleChain)      // GET
	http.HandleFunc("/wallet/new", handleNewWallet) // GET
	http.HandleFunc("/tx/new", handleNewTx)     // POST with tx JSON, header X-Pubkey
	http.HandleFunc("/tx/pool", handleTxPool)   // GET

	log.Println("PURI-CHAIN node running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
