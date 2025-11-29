package main

import (
	"log"
	"math/big"
	"crypto/elliptic"
)

func main() {
	InitChain()
	// Set difficulty or leave default
	Difficulty = 4

	// start server
	go func() {
		startServer()
	}()

	// simple miner start message
	log.Println("PURI-CHAIN node started. Mine with: curl 'http://localhost:8080/mine?miner=YOUR_ADDRESS'")
	select {} // block forever
}
