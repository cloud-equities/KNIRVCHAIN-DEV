package main

import (
	"KNIRVCHAIN-DEV/constants"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type SmartContract struct {
	Code []byte
	Data []byte
}
type Transaction struct {
	From            string `json:"from"`
	To              string `json:"to"`
	Value           uint64 `json:"value"`
	Timestamp       int64  `json:"timestamp"`
	Status          string `json:"status"`
	TransactionHash string `json:"transaction_hash"`
	PublicKey       string `json:"public_key"`
	Signature       string `json:"signature"`

	Data []byte `json:"data"`
}

func (tx *Transaction) VerifyTxn() bool { // to be implemented at later stages when digital signature methods are fully implemented for validating.
	return true // just do basic implementation that every signature passes as an valid method call and as a verification for transactions
}

type Blockchain struct { // Removed this old and wrong implementation struct
}

type ConsensusManager struct { // Implement consensus lock functionality.
	miningLock     bool
	syncState      chan bool
	longestChain   string
	updateRequired bool
	mu             sync.Mutex
}

func NewConsensusManager() *ConsensusManager {
	return &ConsensusManager{
		miningLock:     false,           // do not lock by default.
		syncState:      make(chan bool), // Use data channels for state to signal events or blocking or resource usage or synchronizing operations for resources by state data information during method implementations of goroutines
		updateRequired: false,           // make false, to not run updates until first process requires and sets it to true to correctly make chain use safe to prevent possible corrupted state and logic from method implementations.
	}
}
func (cm *ConsensusManager) lockMining() { // if there are concurrency related issues while creating the genesis block, while the code attempts to make calculations during initial start, prevent that from happening by implementing lock and unlock functionality during object creation to ensure program correctness.
	cm.mu.Lock()
	defer cm.mu.Unlock() // Always use mutex unlock
	cm.miningLock = true // lock state set to prevent race condition implementation bugs
}

func (cm *ConsensusManager) unlockMining() { //  prevent potential implementation state problems, so unlock whenever resource does not requires lock, as those locks or states should only apply only when absolutely necessary during specific cases of your code execution workflows where they are implemented, while testing those logic for all execution and not by a some hardcoded approach when developing a software with such requirements.

	cm.mu.Lock()

	defer cm.mu.Unlock()  // Always unlock your method implementations
	cm.miningLock = false //unlock resource so logic continues
}
func (cm *ConsensusManager) getMiningLockState() bool { // create a simple getter for resource states.

	cm.mu.Lock()         //lock value access for concurrency safety when using it by goroutines
	defer cm.mu.Unlock() // release value after read has been complete.
	return cm.miningLock // avoid race conditions and always ensure the method or the data/variable for testing or state checking access, has a properly set resource lock state or implementation.
}
func (cm *ConsensusManager) getSyncState() chan bool { // data sync mechanism to handle goroutines synchronization of access to resources shared between implementations or components with object or structs or values that are mutable
	return cm.syncState // if a pointer of an object is being accessed from other implementations methods using goroutines during software test/workflow logic and validations.
}
func (cm *ConsensusManager) setLongestChain(hash string) { // update state value when it needs it during testing of correct logic with implementations
	cm.mu.Lock()           // lock for writing safety on object state access when used with methods during concurrent operation, if goroutine related usage to avoid non deterministic system behavior that would cause implementation related bugs in object logic for specific implementations tests validations that would fail for incorrect state changes or value implementations.
	defer cm.mu.Unlock()   // Always use `Unlock`, once process complete.
	cm.longestChain = hash // update value to indicate which chain or what value the state implementation should use or apply as correct during its method/testing implementations of the project being validated in testing or for standard or production logic when creating implementations of code methods.
}

func (cm *ConsensusManager) setUpdateRequired(val bool) { // update state to perform or not perform particular code using a value being passed to implement different workflow when method logic is being executed in current execution context, which is why we require synchronization, with data primitives and struct locks if their method implementations must also perform validation or state checks.
	cm.mu.Lock()            // create lock to ensure correct update during the concurrency or when multiple components are calling, or reading and writing during program/test logic when object are under a concurrency context using goroutines, and also if it applies that multiple processes can or will access at specific times, for safe usage without corrupted state issues for variables implementation details or from method implementations logic behavior with type parameters using or modifying such state.
	defer cm.mu.Unlock()    // unlock whenever complete for resources that use them or set their state that use mutex.
	cm.updateRequired = val // apply changes during a call to these methods.
}
func (cm *ConsensusManager) getUpdateRequired() bool { // create method that reads object state
	cm.mu.Lock()             //create a local object lock for a given variable during state check that must only execute safely, on given time, if is a pointer variable used or shared during program runtime execution during normal workflow state, using goroutine or if it could cause concurrency problems when using that data directly.
	defer cm.mu.Unlock()     // Use `defer` unlock methods to prevent blocking code or resources for too long after implementation steps have been performed by the methods for those state/variable checks from the object and it's workflow execution process.
	return cm.updateRequired // returns object's update state check based on the latest known or current state based on what implementations set up in code logic.
}

func NewTransaction(from, to string, value uint64, data []byte) *Transaction {
	transaction := new(Transaction)

	transaction.From = from
	transaction.To = to
	transaction.Value = value
	transaction.Timestamp = time.Now().Unix()
	transaction.Data = data
	transaction.Status = constants.PENDING // for now by default
	transaction.TransactionHash = generateHash([]byte{}, transaction.Timestamp, fmt.Sprintf("%s%s%d", from, to, value), 0)

	return transaction

}
func generateHash(prevHash []byte, timestamp int64, data string, nonce int) string {

	dataToHash := fmt.Sprintf("%s%d%s%d", string(prevHash), timestamp, data, nonce) // remove as it is incorrect use of Hash on our Block
	h := sha256.Sum256([]byte(dataToHash))
	return hex.EncodeToString(h[:])
}
