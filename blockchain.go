package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Block struct {
	Hash      string `json:"hash"`
	PrevHash  string `json:"prev_hash"`
	Timestamp int64  `json:"timestamp"`
	Data      string `json:"data"`
	Nonce     int    `json:"nonce"`
}
type Blockchain struct {
	Chain        []*Block        `json:"chain"`
	ChainAddress string          `json:"chain_address"`
	Peers        map[string]bool `json:"peers"`
	mu           sync.Mutex
	db           *DBClient
}

func NewBlockchain(genesisBlock *Block, chainAddress string, db *DBClient) *Blockchain {
	bc := &Blockchain{
		Chain:        []*Block{genesisBlock},
		ChainAddress: chainAddress,
		Peers:        map[string]bool{},
		db:           db,
	}
	if db != nil {
		if lastBlock, err := db.GetLastBlock(); err == nil {
			if lastBlock != nil {
				bc.Chain = []*Block{lastBlock}
			}
		} else {
			log.Println("could not retrieve block from database", err)
		}
	}
	return bc
}

func (bc *Blockchain) CreateBlock(data string, prevHash string, nonce int) *Block {
	timestamp := time.Now().Unix()
	block := &Block{
		Hash:      generateHash(prevHash, timestamp, data, nonce),
		PrevHash:  prevHash,
		Timestamp: timestamp,
		Data:      data,
		Nonce:     nonce,
	}
	return block

}

func (bc *Blockchain) AddBlock(block *Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	if len(bc.Chain) == 0 {
		bc.Chain = append(bc.Chain, block)
		if bc.db != nil {
			return bc.db.InsertBlock(block)
		}
		return nil
	}
	lastBlock := bc.Chain[len(bc.Chain)-1]
	if block.PrevHash != lastBlock.Hash {
		return fmt.Errorf("invalid previous hash for the block")
	}
	bc.Chain = append(bc.Chain, block)
	if bc.db != nil {
		return bc.db.InsertBlock(block)
	}
	return nil

}
func (bc *Blockchain) ProofOfWorkMining(minerAddress string) {
	fmt.Println("Starting the mining process...")
	for {
		lastBlock := bc.GetLastBlock()
		nonce := 0
		if lastBlock == nil {
			genesisBlock := &Block{
				Hash:      "genesis_block_hash",
				PrevHash:  "",
				Timestamp: time.Now().Unix(),
				Data:      "genesis_block_data",
				Nonce:     0,
			}
			bc.AddBlock(genesisBlock)
			fmt.Println("Genesis block mined.")
			continue
		}
		for {
			hash := generateHash(lastBlock.Hash, time.Now().Unix(), "transaction", nonce) // replace hardcoded data with transactions
			if strings.HasPrefix(hash, "000") {                                           // proof of work prefix check
				block := bc.CreateBlock("transaction", lastBlock.Hash, nonce)
				if err := bc.AddBlock(block); err != nil {
					log.Println("error adding mined block:", err)
				}
				fmt.Println("Mined block number:", len(bc.Chain), "by", minerAddress)
				break
			}
			nonce++
		}
		time.Sleep(5 * time.Second)
	}
}
func (bc *Blockchain) GetLastBlock() *Block {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	if len(bc.Chain) == 0 {
		if bc.db != nil {
			lastBlock, err := bc.db.GetLastBlock()
			if err != nil {
				log.Println("failed to get last block from db", err)
			}
			return lastBlock
		}
		return nil
	}
	return bc.Chain[len(bc.Chain)-1]
}
func generateHash(prevHash string, timestamp int64, data string, nonce int) string {
	dataToHash := fmt.Sprintf("%s%d%s%d", prevHash, timestamp, data, nonce)
	hash := sha256.Sum256([]byte(dataToHash))
	return hex.EncodeToString(hash[:])
}
func (bc *Blockchain) DialAndUpdatePeers() {

	for {
		for peer := range bc.Peers {
			if peer != bc.ChainAddress {
				go func(peer string) {
					bc.updateChainFromPeer(peer)
				}(peer)
			}
		}
		time.Sleep(10 * time.Second)
	}
}
func (bc *Blockchain) updateChainFromPeer(peer string) {
	url := peer + "/chain"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching chain from peer %s: %v", peer, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Peer %s returned status %d", peer, resp.StatusCode)
		return
	}
	var peerChain []*Block
	if err := json.NewDecoder(resp.Body).Decode(&peerChain); err != nil {
		log.Printf("Error decoding chain from peer %s: %v", peer, err)
		return
	}

	bc.mu.Lock()
	defer bc.mu.Unlock()
	if len(peerChain) > len(bc.Chain) {
		bc.Chain = peerChain
		if bc.db != nil {
			for _, block := range peerChain {
				if err := bc.db.InsertBlock(block); err != nil {
					log.Println("failed to insert block when syncing blockchain", err)
				}
			}
		}
		log.Printf("Chain updated from peer: %s\n", peer)
	}
}
func (bc *Blockchain) RunConsensus() {
	fmt.Println("Starting the consensus algorithm...")
	for {
		for peer := range bc.Peers {
			if peer != bc.ChainAddress {
				go func(peer string) {
					bc.updateChainFromPeer(peer)
				}(peer)
			}
		}
		time.Sleep(15 * time.Second)
	}
}
