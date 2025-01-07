package main

import (
	"KNIRVCHAIN-DEV/constants"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"time"
)

type ConsensusManager struct { // Implement consensus lock functionality.
	miningLock      bool
	syncState       chan bool
	longestChain    []*Block
	updateRequired  bool
	Blockchain      *BlockchainStruct
	MiningLocked    bool
	TransactionPool *TransactionPool

	mu sync.Mutex
}

type TransactionPool struct {
	transactionPool []*Transaction
	mu              sync.Mutex
}

func (tp *TransactionPool) Lock() {
	tp.mu.Lock()
}

func (tp *TransactionPool) Unlock() {
	tp.mu.Unlock()
}

// Get returns a copy of the transaction pool slice. This is important to avoid race
// conditions when reading the pool concurrently.
func (tp *TransactionPool) Get() []*Transaction {
	tp.Lock()
	defer tp.Unlock()
	cpy := make([]*Transaction, len(tp.transactionPool))
	copy(cpy, tp.transactionPool) // Create a copy of the underlying slice to return
	return cpy
}

// Set updates the transaction pool. It is important to always use Set() when
// changing the transaction pool to avoid issues related to concurrency.
func (tp *TransactionPool) Set(txns []*Transaction) {
	tp.Lock()
	defer tp.Unlock()
	tp.transactionPool = txns
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		transactionPool: make([]*Transaction, 0), // Initialize the pool
	}
}

func NewConsensusManager(blockchain *BlockchainStruct, transactionPool *TransactionPool) *ConsensusManager { // Created an new struct type.

	return &ConsensusManager{
		miningLock:      false,
		syncState:       make(chan bool),
		updateRequired:  false,
		MiningLocked:    false,
		Blockchain:      blockchain,
		TransactionPool: transactionPool, // setting the blockchain struct type to be used in the consensus manager struct type for the object state to be used in the project when it is being created and used in the project
	}

}

func (cm *ConsensusManager) lockMining() { // lock mining

	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.miningLock = true

}

func (cm *ConsensusManager) unlockMining() { // lock implementation using workflows from struct implementations safely where parameter requires that code.

	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.miningLock = false
}
func (cm *ConsensusManager) getMiningLockState() bool { // access those types safely with a data object by getting value with correct types for local resource to avoid race conditions when methods are being validated with state using that workflow (where type properties where modified)
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.miningLock
}
func (cm *ConsensusManager) getSyncState() chan bool { // data channel state to also guarantee workflows validations based on concurrency

	return cm.syncState // always access object workflow with interfaces types of validations (implement a standard that they have properties access for memory/object local state usage safety in project to be portable during implementations where tests using method must have safe resources that must match also code parameter/type validations when these are shared and used).

}
func (cm *ConsensusManager) setLongestChain(longestChain []*Block) {
	cm.mu.Lock()         // correct locking with parameter validation (all local variables that modify object states must have protections via proper access during tests that methods also requires such states be validated via object types )
	defer cm.mu.Unlock() // correct way using the interface or methods workflow that is required by struct to perform such operation as part of validation implementation.

	cm.longestChain = longestChain // safely implements update during a local state access in type parameters from object instances with validation or system requirement (memory is protected if it uses shared workflows of software with mutex using validations on workflows).

}
func (cm *ConsensusManager) setUpdateRequired(val bool) {
	cm.mu.Lock() //lock
	defer cm.mu.Unlock()
	cm.updateRequired = val // use implementation requirements for setting data or object type state properties workflow with type validations in memory when project also relies on it to also test/validation and implementations methods during data validation for project (such operations can't corrupt data when is being accessed in code)
}
func (cm *ConsensusManager) getUpdateRequired() bool { // Get local state that method use during test implementation workflow when workflows validates types via interfaces or method parameters of struct data using the validation implementations.

	defer cm.mu.Unlock()

	return cm.updateRequired

}

func FetchLastNBlocks(peer string) (*BlockchainStruct, error) {
	ourURL := fmt.Sprintf("%s/", peer)
	resp, err := http.Get(ourURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//unmarshal the response
	bs := new(BlockchainStruct)
	return bs, json.NewDecoder(resp.Body).Decode(bs)
}

func (cm *ConsensusManager) RunConsensus() {
	cm.MiningLocked = false
	for {
		log.Println("Starting the consensus algorithm...")

		// check if the longest chain is our
		if cm.Blockchain == nil { // Check for nil Blockchain
			log.Println("Blockchain is nil. Skipping consensus iteration.")
			time.Sleep(constants.CONSENSUS_PAUSE_TIME * time.Second)
			continue
		}

		if cm.Blockchain.Peers == nil { // Check if Peers is nil
			log.Println("Peers list is nil. Skipping peer synchronization.")
			time.Sleep(constants.CONSENSUS_PAUSE_TIME * time.Second)
			continue // Skip this iteration
		}

		longestChain := cm.Blockchain.Blocks

		longestChainIsOur := true
		for peer, status := range cm.Blockchain.Peers {
			if status {
				bc1, err := FetchLastNBlocks(peer)
				if err != nil {
					log.Println("Error while  fetching last n blocks from peer:", peer, "Error:", err.Error())
					continue
				}

				if len(bc1.Blocks) > len(longestChain) {
					longestChain = bc1.Blocks

					longestChainIsOur = false
				}
			}
		}

		if longestChainIsOur {
			log.Println("My chain is longest, thus I am not updating my blockchain")
			time.Sleep(constants.CONSENSUS_PAUSE_TIME * time.Second)
			continue
		}

		if verifyLastNBlocks(longestChain) {
			// stop the Mining until updation
			cm.MiningLocked = true
			cm.UpdateBlockchain(longestChain)
			// restart the Mining as updation is complete
			cm.MiningLocked = false
			log.Println("Updation of Blockchain complete !!!")
		} else {
			log.Println("Chain Verification Failed, Hence not updating my blockchain")
		}
		time.Sleep(constants.CONSENSUS_PAUSE_TIME * time.Second)
	}

}

func verifyLastNBlocks(_ []*Block) bool {
	// Implement your verification logic here
	return true // Replace with actual verification
}

func (cm *ConsensusManager) UpdateBlockchain(longestChain []*Block) {
	// Lock before updating blockchain to ensure thread safety
	cm.Blockchain.Lock()
	defer cm.Blockchain.Unlock()

	// Step 1: Validate the integrity of the incoming longestChain
	if !cm.validateLongestChain(longestChain) {
		log.Println("Validation failed: longest chain is not valid.")
		return
	}

	// Step 2: Deep copy the longestChain to avoid unintended modifications
	validatedChain := make([]*Block, len(longestChain))
	for i, block := range longestChain {
		validatedChain[i] = block.DeepCopy() // Assuming Block has a DeepCopy method
	}

	// Step 3: Update the blockchain with the validated longestChain
	cm.Blockchain.SetBlocks(validatedChain)

	// Step 4: Update other related data structures, if necessary
	cm.updateTransactionPool(validatedChain)

	// Step 5: Persist the updated blockchain to the database
	db := NewLevelDB() // Create a new instance of the database
	defer db.Close()   // Ensure the database is closed properly
	err := cm.Blockchain.PutIntoDb(db, cm.Blockchain.ChainAddress)
	if err != nil {
		log.Printf("Error updating blockchain to DB: %v", err)
		return
	}

	log.Println("Blockchain successfully updated.")
}

// validateLongestChain checks the validity of the provided chain
func (cm *ConsensusManager) validateLongestChain(chain []*Block) bool {
	if len(chain) == 0 {
		return false
	}

	for i := 1; i < len(chain); i++ {
		// Ensure each block is linked correctly
		if !bytes.Equal(chain[i].PrevHash, chain[i-1].Hash()) {
			return false
		}
		if chain[i].Transactions == nil {
			log.Println("Transactions are nil in block:", i)
			return false
		}

		// Additional validations, e.g., block difficulty, timestamps, etc.
		if !chain[i].IsValid() {
			return false
		}
	}

	return true
}

// updateTransactionPool removes confirmed transactions from the transaction pool
func (cm *ConsensusManager) updateTransactionPool(validatedChain []*Block) {
	if cm.TransactionPool == nil { // Handle nil TransactionPool
		log.Println("Transaction pool is nil. Skipping update.")
		return
	}

	cm.TransactionPool.Lock()
	defer cm.TransactionPool.Unlock()

	// Create a map of confirmed transactions

	confirmedTxs := make(map[string]bool)
	for _, block := range validatedChain {
		for _, tx := range block.Transactions {
			confirmedTxs[tx.PublicKey] = true
		}
	}

	var unconfirmedTxs []*Transaction
	for _, tx := range cm.TransactionPool.Get() {
		if _, ok := confirmedTxs[tx.TransactionHash]; !ok {
			unconfirmedTxs = append(unconfirmedTxs, tx)
		}
	}
	cm.TransactionPool.Set(unconfirmedTxs)
}
