package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type BlockchainServer struct {
	port          uint64
	BlockchainPtr *Blockchain
	chainAddress  string
}

func NewBlockchainServer(port uint64, blockchain *Blockchain, chainAddress string) *BlockchainServer {
	return &BlockchainServer{
		port:          port,
		BlockchainPtr: blockchain,
		chainAddress:  chainAddress,
	}
}

func (bcs *BlockchainServer) Start() {
	http.HandleFunc("/chain", bcs.handleGetChain)
	http.HandleFunc("/peers", bcs.handleGetPeers)
	http.HandleFunc("/add_peer", bcs.handleAddPeer)

	log.Println("Starting server on port: ", bcs.port)
	if err := http.ListenAndServe(":"+strconv.Itoa(int(bcs.port)), nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}
func (bcs *BlockchainServer) handleGetChain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	chain := bcs.BlockchainPtr.Chain

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(chain); err != nil {
		log.Println("failed to encode response", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func (bcs *BlockchainServer) handleGetPeers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	peers := bcs.BlockchainPtr.Peers
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(peers); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
func (bcs *BlockchainServer) handleAddPeer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var newPeer string
	if err := json.NewDecoder(r.Body).Decode(&newPeer); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	bcs.BlockchainPtr.Peers[newPeer] = true
	fmt.Fprintf(w, "Peer %s added successfully\n", newPeer)
}
