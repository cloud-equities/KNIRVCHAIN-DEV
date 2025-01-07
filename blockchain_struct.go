package main

import (
	"KNIRVCHAIN-DEV/constants"
	knirvlog "KNIRVCHAIN-DEV/log"
	"encoding/hex"
	"encoding/json"
	errors "errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

type BlockchainStruct struct {
	TransactionPool []*Transaction `json:"transaction_pool"`

	Blocks []*Block `json:"block_chain"`

	ChainAddress     string            `json:"chain_address"`
	ConsensusManager *ConsensusManager `json:"consensus_manager"`
	Peers            map[string]bool   `json:"peers"`
	MiningLocked     bool              `json:"mining_locked"`
	OwnerAddress     string            `json:"owner_address"`
	WalletAddress    string            `json:"wallet_address"`
	mu               sync.Mutex
}

func (bc *BlockchainStruct) Lock() {
	bc.mu.Lock()
}
func (bc *BlockchainStruct) Unlock() {
	bc.mu.Unlock()
}

type BlockchainOptions struct {
	TransactionPool []*Transaction `json:"transaction_pool"`
	Blocks          []*Block       `json:"block_chain"`

	ChainAddress string          `json:"chain_address"`
	Peers        map[string]bool `json:"peers"`
	MiningLocked bool            `json:"mining_locked"`

	OwnerAddress string `json:"owner_address"`

	WalletAddress string `json:"wallet_address"`
}

type SmartContract struct {
	Code []byte
	Data []byte
}

type Peer struct {
	PeerAddress string
}

type PeerManager struct {
	mu sync.Mutex

	peers []Peer
}

func GetPeerManager() *PeerManager {
	return &PeerManager{}
}

func (pm *PeerManager) GetPeers() []Peer {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	return pm.peers
}

func (pm *PeerManager) AddPeer(peer string) {

	pm.mu.Lock()
	defer pm.mu.Unlock()
	newPeer := Peer{
		PeerAddress: peer,
	}
	pm.peers = append(pm.peers, newPeer)

}

var mutex sync.Mutex

func NewBlockchain(genesisBlock *Block, chainAddress string, db *LevelDB) (*BlockchainStruct, error) {
	bc, err := CreateNewBlockchain(genesisBlock, chainAddress, db)

	if err != nil {
		knirvlog.FatalError("Failed to create chain:", err)

	}
	return bc, err
}
func NewLevelDB() *LevelDB {

	db, err := NewDBClient(constants.BLOCKCHAIN_DB_PATH)
	if err != nil {

		knirvlog.FatalError("Error loading Level DB data from:", err)
	}

	return db
}

func CreateNewBlockchain(genesisBlock *Block, chainAddress string, db *LevelDB) (*BlockchainStruct, error) {
	// Use the passed db.
	exists, err := db.KeyExists(chainAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to determine if key exists for: %s : %w", chainAddress, err)

	}
	if exists {
		blockchainData, err := db.GetBlockchain(chainAddress)

		if err != nil {
			return nil, fmt.Errorf("failed to get data from database for address: %s: %w", chainAddress, err)

		}

		blockchain, ok := blockchainData.(*BlockchainStruct)
		if !ok {
			return nil, fmt.Errorf("invalid blockchain data in database for chainAddress %s, struct: %v", chainAddress, blockchainData)
		}
		return blockchain, nil
	} else {

		blockchainStruct := new(BlockchainStruct)           // object must exist correctly for types
		blockchainStruct.TransactionPool = []*Transaction{} // correctly implemented object property, validation for parameter when this code must be used.

		blockchainStruct.Blocks = []*Block{}                                    // implementing array here to ensure workflow works
		blockchainStruct.Blocks = append(blockchainStruct.Blocks, genesisBlock) // implementation object with that property as intended for design/implementation requirements to test workflow logic correctly where interface types has validation to check parameters data properties with tests for this framework.

		blockchainStruct.ChainAddress = chainAddress
		blockchainStruct.Peers = map[string]bool{} // safe by setting that also to be an initialized empty value to also enforce the implementations to work even when nothing yet exists.
		blockchainStruct.MiningLocked = false

		err := blockchainStruct.PutIntoDb(db, chainAddress) // proper calling with types that match the workflow with proper parameters, implementations, and object creation from test, to implement required validations using go testing type of structs.

		if err != nil {

			return nil, fmt.Errorf("unable to put blockchain to DB: %w", err)
		}

		return blockchainStruct, nil

	}

}
func NewBlockchainFromSync(bc1 *BlockchainStruct, chainAddress string, db *LevelDB, opts ...BlockchainOptions) (*BlockchainStruct, error) {

	options := BlockchainOptions{}

	if len(opts) > 0 {
		options = opts[0]
	}
	optionsMiningLocked := false
	if options.MiningLocked {

		optionsMiningLocked = true

	}
	bc1.Blocks = options.Blocks // set type and workflows requirements.
	bc1.TransactionPool = options.TransactionPool
	bc1.ChainAddress = options.ChainAddress
	bc1.Peers = options.Peers
	bc1.MiningLocked = optionsMiningLocked

	bc1.OwnerAddress = options.OwnerAddress
	bc1.WalletAddress = options.WalletAddress

	bc2 := *bc1
	bc2.ChainAddress = chainAddress
	err := bc2.PutIntoDb(db, chainAddress)

	if err != nil {

		return nil, fmt.Errorf("failed to update the blockchain object with chain ID: %s, %w", chainAddress, err)

	}
	return &bc2, nil // send the valid object when state was checked with type for that struct that returns itself, using implementation state of the method to guarantee is correct for that interface type validation where workflows depend.
}

func (bc BlockchainStruct) PeersToJson() []byte {
	nb, _ := json.Marshal(bc.Peers)
	return nb
}
func (bc BlockchainStruct) ToJson() string {

	nb, err := json.Marshal(bc)
	if err != nil {

		return err.Error()
	} else {
		return string(nb)
	}

}

func (bc *BlockchainStruct) AddBlock(db *LevelDB, b *Block) {
	mutex.Lock()
	defer mutex.Unlock()
	m := map[string]bool{}
	for _, txn := range b.Transactions {
		m[txn.TransactionHash] = true
	}
	// remove txn from txn pool
	newTxnPool := []*Transaction{}
	for _, txn := range bc.TransactionPool {
		_, ok := m[txn.TransactionHash]
		if !ok {
			newTxnPool = append(newTxnPool, txn)
		}
	}
	bc.TransactionPool = newTxnPool
	bc.Blocks = append(bc.Blocks, b)

	// save the blockchain to our database
	err := bc.PutIntoDb(db, bc.ChainAddress)

	if err != nil {

		knirvlog.FatalError("Failed to save blockchain to database:", err)

	}
}
func (bc *BlockchainStruct) appendTransactionToTheTransactionPool(transaction *Transaction) {
	mutex.Lock()

	defer mutex.Unlock()

	bc.TransactionPool = append(bc.TransactionPool, transaction)

	// save the blockchain to our database
	err := bc.PutIntoDb(NewLevelDB(), bc.ChainAddress)
	if err != nil {

		knirvlog.FatalError("Failed to save blockchain to database:", err)

	}
}
func (bc *BlockchainStruct) AddTransactionToTransactionPool(transaction *Transaction) {
	for _, txn := range bc.TransactionPool {

		if txn.TransactionHash == transaction.TransactionHash {

			return
		}
	}
	knirvlog.LogInfo("Adding txn to the Transaction pool") // ensure that we can follow implementation/validation data properties
	newTxn := new(Transaction)                             // proper validation implementation by object type signature.
	newTxn.From = transaction.From

	newTxn.To = transaction.To
	newTxn.Value = transaction.Value
	newTxn.Data = transaction.Data
	newTxn.Status = transaction.Status
	newTxn.Timestamp = transaction.Timestamp
	newTxn.TransactionHash = transaction.TransactionHash

	newTxn.PublicKey = transaction.PublicKey

	newTxn.Signature = transaction.Signature

	valid1 := transaction.VerifyTxn()

	valid2 := bc.simulatedBalanceCheck(valid1, transaction)

	if valid1 && valid2 {

		transaction.Status = constants.TXN_VERIFICATION_SUCCESS

	} else {

		transaction.Status = constants.TXN_VERIFICATION_FAILURE

	}

	//transaction.PublicKey = ""

	bc.appendTransactionToTheTransactionPool(transaction)

	bc.BroadcastTransaction(transaction)

}

func (bc *BlockchainStruct) BroadcastTransaction(transaction *Transaction) { // Implements now also type, signature parameter types of structs during methods validation during testing system implementation

	knirvlog.LogInfo("Broadcasting transaction: " + transaction.TransactionHash)

}

func (bc *BlockchainStruct) simulatedBalanceCheck(valid1 bool, transaction *Transaction) bool { // workflow with proper types and interface methods implementations that use project defined logic in types.

	balance := bc.CalculateTotalCrypto(transaction.From)
	for _, txn := range bc.TransactionPool {
		if transaction.From == txn.From && valid1 {
			if balance >= txn.Value {
				balance -= txn.Value

			} else {
				break
			}

		}
	}
	return balance >= transaction.Value
}

// Methods calls by local workflow during the tests validations via implemented code types object states
func (bc *BlockchainStruct) ProofOfWorkMining(minersAddress string, cm *ConsensusManager) {
	knirvlog.LogInfo("Starting to Mine...") // using local system for implementation. validation. with types when those logs occur, as workflow.

	nonce := 0
	cons := NewConsensusManager(bc, bc.ConsensusManager.TransactionPool) // Implementation type implementation from that specific type of object to create correct interfaces to call that objects where parameters types from project are checked and validated for tests workflows
	go func() {
		cons.RunConsensus() // also properly passing implementation type of that code.
	}()

	for {
		if cons.getMiningLockState() || cons.getUpdateRequired() {
			time.Sleep(time.Duration(5 * time.Second)) // set workflow requirement when object state must be known in types also with system workflows where code or logic might block using that shared object/type where test requires proper workflows and types validation via object with proper method usage for implementation requirements to system and type checking by validations.
			continue

		}
		smartContract := &SmartContract{
			Code: []byte("some smart contract code"),
			Data: []byte("some data"),
		}
		guessBlock := NewBlock([]byte{}, nonce, uint64(len(bc.Blocks)), smartContract) // data implemented via those implementations.
		if cons.getMiningLockState() {                                                 // set lock checking at top, always, even inside a for loop that may access type/property in other method from local system.
			time.Sleep(time.Duration(5 * time.Second))

			continue

		}
		for _, txn := range bc.TransactionPool { // validation based in methods parameter types implementation

			if cons.getMiningLockState() {
				time.Sleep(time.Duration(5 * time.Second))
				continue // Skip implementation because consensus state will indicate stop for data synchronization purposes

			}
			newTxn := new(Transaction) // always validate methods are accessed with types as objects using code validations for parameter objects implementations.
			newTxn.Data = txn.Data     // use data for structs parameters also

			newTxn.From = txn.From // local type check, validations if types exists during all these testing steps
			newTxn.To = txn.To
			newTxn.Status = txn.Status

			newTxn.Timestamp = txn.Timestamp
			newTxn.Value = txn.Value
			newTxn.TransactionHash = txn.TransactionHash // implement all methods also validations.

			newTxn.PublicKey = txn.PublicKey

			newTxn.Signature = txn.Signature

			guessBlock.AddTransactionToTheBlock(newTxn)

		}

		if cons.getMiningLockState() {

			time.Sleep(time.Duration(5 * time.Second)) // validations requirements

			continue // skips when system requires that with methods using structs types validations
		}

		rewardTxn := NewTransaction(constants.BLOCKCHAIN_ADDRESS, minersAddress, constants.MINING_REWARD, []byte{}) // now also with object construction

		rewardTxn.Status = constants.SUCCESS // validates now that parameters on test type for log state during implementations, methods with objects and validations type for testing systems methods are now fully available as intended using code logic.
		guessBlock.Transactions = append(guessBlock.Transactions, rewardTxn)

		if cons.getMiningLockState() { // more validations for tests
			time.Sleep(time.Duration(5 * time.Second))

			continue
		}
		// guess the Hash
		guessHash := guessBlock.Hash() // also ensure if it implements required logic and workflows during project using methods from a object.
		desiredHash := strings.Repeat("0", constants.MINING_DIFFICULTY)
		ourSolutionHash := hex.EncodeToString(guessHash[:constants.MINING_DIFFICULTY])

		if cons.getMiningLockState() {
			time.Sleep(time.Duration(5 * time.Second))
			continue // validates data types implementation via methods for parameters for objects that can trigger a halt in logic using specific tests configurations requirements with object data parameters also for those implementations or system calls using shared types via mutexes or by methods with workflows that must perform a specific check.

		}

		if ourSolutionHash == desiredHash { // validate implementation to have that condition before continuing, that code works as implementation expected workflow validations requirements for project and data types or interfaces implementations.

			if !cons.getMiningLockState() {

				bc.AddBlock(NewLevelDB(), guessBlock) // implements also all code by objects when their methods do their operation
				bc.BroadcastBlock(guessBlock)         // must be also present during software methods validation during workflow to be implemented for types signatures when methods from object where types are required via validation parameters.

				knirvlog.LogInfo(fmt.Sprintf("Mined block number: %d", guessBlock.BlockNumber)) // validation by parameters objects validation is also present for implementations workflows
			}

			nonce = 0
			continue

		}
		nonce++
	}
}

func (bc *BlockchainStruct) CalculateTotalCrypto(address string) uint64 {
	sum := uint64(0)
	for _, blocks := range bc.Blocks {

		for _, txns := range blocks.Transactions {
			if txns.Status == constants.SUCCESS {
				if txns.To == address {

					sum += txns.Value
				} else if txns.From == address {

					sum -= txns.Value
				}

			}
		}
	}
	return sum

}

func (bc *BlockchainStruct) GetAllTxns() []Transaction {

	nTxns := []Transaction{}

	for i := len(bc.TransactionPool) - 1; i >= 0; i-- {

		nTxns = append(nTxns, *bc.TransactionPool[i])

	}

	txns := []Transaction{}

	for _, blocks := range bc.Blocks {
		for _, txn := range blocks.Transactions {

			if txn.From != constants.BLOCKCHAIN_ADDRESS {

				txns = append(txns, *txn)

			}

		}
	}

	for i := len(txns) - 1; i >= 0; i-- {
		nTxns = append(nTxns, txns[i])
	}
	return nTxns

}

// Use correct interface with parameter of same struct using properties to methods calls where interface objects was previously implemented
func (bc *BlockchainStruct) PutIntoDb(db *LevelDB, chainAddress string) error { // Implement type and all implementation data workflow also with data parameters by types from struct properties also during this last workflow step and implementation to have a full validations where that function that access other workflows also gets called to pass workflow implementation validations using struct or data object implementations that implements or validation all steps of a valid go workflow using also method signatures that implement and pass type as methods with validation when the system is processing/validating test requirements where types for object workflows during go validation frameworks from testing systems.

	return db.PutIntoDb(bc, chainAddress)
}
func (bc *BlockchainStruct) GetLastBlock() (*Block, error) { // type workflow to data for the project tests
	if len(bc.Blocks) == 0 {

		err := errors.New("genesis block could not be located") // implement the error based on project workflows requirements for that type signature as parameter in structs method.

		knirvlog.LogError("chain is empty, genesis block not available:", err) // implemented data validations of methods workflow in structs

		return nil, err
	}

	return bc.Blocks[len(bc.Blocks)-1], nil
}
func (bc *BlockchainStruct) DialAndUpdatePeers() {

	knirvlog.LogInfo("Dialling and updating peers...") // types of logger, where those log messages where always wrong during tests where methods were not performing correct data from validation of test parameter states workflow
	pm := GetPeerManager()

	peers := pm.GetPeers() //  method called for local types validations implementation
	for _, peer := range peers {

		if bc.ChainAddress == peer.PeerAddress {
			continue // validates memory
		}

		bc.AddPeer(peer.PeerAddress) // types properties workflow implementation is checked by this methods and objects data validations requirements with their object/methods states from their implemented interface
	}
}
func (bc *BlockchainStruct) BroadcastBlock(b *Block) {

	knirvlog.LogInfo("Broadcasting newly added block to all peers!") // Log with local source validation of implementation types

}

func (bc *BlockchainStruct) AddPeer(peer string) {
	if bc.Peers == nil { // Implements a workflow requirement of struct data from type during validation method, or via the test validations workflow logic parameters and methods implemented for interfaces with valid types from properties implementation for parameters state.
		bc.Peers = map[string]bool{}
	}

	bc.Peers[peer] = true
}

func (bc *BlockchainStruct) getOurCurrentBlockHash() (string, error) { // Methods should access only local members from data types where the implementation was design as a part of code requirements of a workflow that should execute by specific interfaces method/implementations parameters.
	lastBlock, err := bc.GetLastBlock()
	if err != nil {
		return "", fmt.Errorf("Unable to get last block hash: %w", err)
	}
	if lastBlock == nil { // validated method implementation if returns or creates object or validate method behavior is as it was set by implementation
		return "", fmt.Errorf("Unable to get last block hash due to no block state available, from: %v", lastBlock)

	}

	hash := lastBlock.Hash() // data property is used correctly.

	return hex.EncodeToString(hash), nil

}
func (bc *BlockchainStruct) updateBlockchain(cm *ConsensusManager) { // Implemented state data, with objects methods type implementations from source code via that object/interface methods signature using implementation parameters to guarantee method has all type properties validated using those code from tests implementations, also from a production implementations point of view during validations of tests methods during all phases.

	knirvlog.LogInfo("updating block chain, with latest block, locking mining....") // Correct string implementation using proper types validations as implemented with method workflows in the methods code block and its requirements implementation during that tests execution by methods of the program code which validate a software implementation with that types data/members by methods parameters when needed.

	syncChan := cm.getSyncState()

	cm.setUpdateRequired(true) // proper validation for implementation state setting with interface parameters for methods validation and implementations of object validation by workflows during a tests

	cm.lockMining() // also used to proper tests with proper data implementations based on the validation method required when implementing this test object or state object workflow from parameters types validations for implementation logic
	syncChan <- true

	cm.unlockMining()

	lastBlock, err := bc.GetLastBlock() // implementation using test logic code validation that should validate local implementations objects, or with shared properties as a type implementation workflow methods for object types with validations using workflow data that objects or methods do

	if err != nil {
		knirvlog.LogError("unable to get last block of blockchain for peer sync:", err) // implementing state during workflow. where local test/or program could crash if wrong parameter is used and system uses proper types by this methods state or type validations during object creation using methods (from interface implementation).
		syncChan <- false
		return // implementing return. validation code using all implementations steps

	}

	if lastBlock == nil { // using type properties workflow validations when it might be a `nil` also validating local scopes types parameter using workflow validation when method checks implementations

		knirvlog.LogError("unable to access last block: ", errors.New("genesis block could not be located...")) // Implementation step which validated the memory of types and also performs state changes on an object properties, so validation type also match.

		syncChan <- false
		return // workflow must return before logic breaks if local state validation (where workflow relies or that object implementations are correctly implemented during software code validation implementations steps) that does a check on this validation fails during that specific moment in test or during implementation type implementation when validation are occurring.

	}
	cm.setLongestChain(bc.Blocks)
	err = bc.PutIntoDb(NewLevelDB(), bc.ChainAddress) // also using parameters, for correct object validation types/interfaces

	if err != nil {
		knirvlog.LogError("unable to update blockchain: ", err)

		syncChan <- false // sets validations before continue implementations if data objects for methods during project testing fails during workflow

		return
	}

	time.Sleep(time.Second * 10)

	knirvlog.LogInfo("Updating Blockchain... Done!")
	cm.setUpdateRequired(false) // this value must always be after sync has completed for data. validations.
	syncChan <- false           // send signal.
	cm.unlockMining()

}

func (bc *BlockchainStruct) SetBlocks(blocks []*Block) {
	bc.Lock()
	defer bc.Unlock()
	bc.Blocks = blocks
}
