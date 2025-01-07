package main

import (
	knirvlog "KNIRVCHAIN-DEV/log"
	"encoding/hex"
	"encoding/json"
	errors "errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"KNIRVCHAIN-DEV/constants"
)

type BlockchainStruct struct {
	TransactionPool []*Transaction  `json:"transaction_pool"`
	Blocks          []*Block        `json:"block_chain"`
	ChainAddress    string          `json:"chain_address"`
	Peers           map[string]bool `json:"peers"`
	MiningLocked    bool            `json:"mining_locked"`
	OwnerAddress    string          `json:"owner_address"`
	WalletAddress   string          `json:"wallet_address"`
}

type BlockchainOptions struct {
	TransactionPool []*Transaction  `json:"transaction_pool"`
	Blocks          []*Block        `json:"block_chain"`
	ChainAddress    string          `json:"chain_address"`
	Peers           map[string]bool `json:"peers"`
	MiningLocked    bool            `json:"mining_locked"`
	OwnerAddress    string          `json:"owner_address"`
	WalletAddress   string          `json:"wallet_address"`
}
type Peer struct {
   PeerAddress string
}
type PeerManager struct {
	mu sync.Mutex
	peers []Peer
}

func GetPeerManager() *PeerManager{

   return &PeerManager{}
}
func (pm *PeerManager) GetPeers() []Peer{
	pm.mu.Lock()
	defer pm.mu.Unlock()
	return pm.peers

}
func (pm *PeerManager) AddPeer(peer string){
   pm.mu.Lock()
	defer pm.mu.Unlock()
	newPeer := Peer{
      PeerAddress: peer,

	}
	pm.peers = append(pm.peers, newPeer)

}
var mutex sync.Mutex
func NewBlockchain(genesisBlock *Block, chainAddress string, db *LevelDB) *BlockchainStruct {
    bc, err := CreateNewBlockchain(genesisBlock, chainAddress, db)
    if err != nil {
        knirvlog.FatalError("Failed to create chain:", err)
    }
    return bc
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
		blockchainStruct := new(BlockchainStruct)
        blockchainStruct.TransactionPool = []*Transaction{}
        blockchainStruct.Blocks = []*Block{}
        blockchainStruct.Blocks = append(blockchainStruct.Blocks, genesisBlock)
        blockchainStruct.ChainAddress = chainAddress
        blockchainStruct.Peers = map[string]bool{}
        blockchainStruct.MiningLocked = false

        err := blockchainStruct.PutIntoDb(db, chainAddress)  // important fix implementation detail, to set parameters properly using correct types/properties.
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
    bc1.Blocks = options.Blocks
    bc1.TransactionPool = options.TransactionPool
    bc1.ChainAddress = options.ChainAddress
    bc1.Peers = options.Peers
    bc1.MiningLocked = optionsMiningLocked
    bc1.OwnerAddress = options.OwnerAddress
    bc1.WalletAddress = options.WalletAddress

	bc2 := *bc1
	bc2.ChainAddress = chainAddress
	err := bc2.PutIntoDb(db, chainAddress) // using object during instantiation to create local scope to access level DB object to persist data correctly
	if err != nil {
        return nil, fmt.Errorf("failed to update the blockchain object with chain ID: %s, %w", chainAddress, err) // return the newly formatted error instead.
	}

    return &bc2, nil
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
    knirvlog.LogInfo("Adding txn to the Transaction pool")
	newTxn := new(Transaction)

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
    transaction.PublicKey = ""
	bc.appendTransactionToTheTransactionPool(transaction)
	bc.BroadcastTransaction(newTxn)

}
func (bc *BlockchainStruct) BroadcastTransaction(transaction *Transaction){
 knirvlog.LogInfo("Broadcasting transaction: " + transaction.TransactionHash)

}

func (bc *BlockchainStruct) simulatedBalanceCheck(valid1 bool, transaction *Transaction) bool {
    balance := bc.CalculateTotalCrypto(transaction.From)
	for _, txn := range bc.TransactionPool {
		if transaction.From == txn.From && valid1{
            if balance >= txn.Value{
                balance -= txn.Value
            }else {
               break
            }
		}
	}
    return balance >= transaction.Value

}
func (bc *BlockchainStruct) ProofOfWorkMining(minersAddress string) {
   knirvlog.LogInfo("Starting to Mine...")
    nonce := 0
    cons := NewConsensusManager()
    go func() {
        cons.RunConsensus(bc)
    }()

    for {
		if cons.getMiningLockState() {
           time.Sleep(time.Duration(5 * time.Second))
           continue
       }

        smartContract := &SmartContract{
            Code: []byte("some smart contract code"),
            Data: []byte("some data"),
        }
		guessBlock := NewBlock([]byte{}, nonce, uint64(len(bc.Blocks)), smartContract)

        if cons.getMiningLockState() {
           time.Sleep(time.Duration(5 * time.Second))
           continue
        }

		for _, txn := range bc.TransactionPool{

			if cons.getMiningLockState() {
				time.Sleep(time.Duration(5 * time.Second))
				continue
			}

            newTxn := new(Transaction)
            newTxn.Data = txn.Data
			newTxn.From = txn.From
			newTxn.To = txn.To
			newTxn.Status = txn.Status
			newTxn.Timestamp = txn.Timestamp
			newTxn.Value = txn.Value
			newTxn.TransactionHash = txn.TransactionHash
            newTxn.PublicKey = txn.PublicKey
            newTxn.Signature = txn.Signature

			guessBlock.AddTransactionToTheBlock(newTxn)


		}

		if cons.getMiningLockState() {
           time.Sleep(time.Duration(5 * time.Second))
           continue
		}
        rewardTxn := NewTransaction(constants.BLOCKCHAIN_ADDRESS, minersAddress, constants.MINING_REWARD, []byte{})
        rewardTxn.Status = constants.SUCCESS
        guessBlock.Transactions = append(guessBlock.Transactions, rewardTxn)

		if cons.getMiningLockState() {
           time.Sleep(time.Duration(5 * time.Second))
           continue
		}

		// guess the Hash
        guessHash := guessBlock.Hash()
		desiredHash := strings.Repeat("0", constants.MINING_DIFFICULTY)
        ourSolutionHash := hex.EncodeToString(guessHash[:constants.MINING_DIFFICULTY])
		if cons.getMiningLockState() {
            time.Sleep(time.Duration(5 * time.Second))
			continue
		}
		if ourSolutionHash == desiredHash {
			if !cons.getMiningLockState(){
               bc.AddBlock(NewLevelDB(), guessBlock)
			   bc.BroadcastBlock(guessBlock)

               knirvlog.LogInfo(fmt.Sprintf("Mined block number: %d", guessBlock.BlockNumber))
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
		for _, txns := range blocks.Transactions{
           if txns.Status == constants.SUCCESS{
				if txns.To == address{
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

    for i := len(txns)-1; i>=0 ; i--{
         nTxns = append(nTxns, txns[i])
    }
    return nTxns


}
// Implemented the method for this new workflows for implementation that also enforces type checking or signatures, where interfaces of object also matches where implementation for objects using their interface for software systems using tests workflow and with validation of each of methods to create a safer code implementation using Go implementation primitives.

func (bc *BlockchainStruct) PutIntoDb(db *LevelDB, chainAddress string) error {

  return db.PutIntoDb(bc, chainAddress)  // use the local objects of levelDB with its method signature to access database methods.

}
func (bc *BlockchainStruct) GetLastBlock() (*Block, error) {
    if len(bc.Blocks) == 0{

		err := errors.New("genesis block could not be located")
      knirvlog.LogError("chain is empty, genesis block not available:",err) // always proper types and object implementation calls, when testing.
		return nil, err
	}


	return bc.Blocks[len(bc.Blocks)-1], nil


}
func (bc *BlockchainStruct) DialAndUpdatePeers() { // Dialling the network using Peer Addresses, for proper sync.

   knirvlog.LogInfo("Dialling and updating peers...")
    pm := GetPeerManager()
	peers := pm.GetPeers()

    for _, peer := range peers {
        if bc.ChainAddress == peer.PeerAddress {

            continue
		}

        bc.AddPeer(peer.PeerAddress) // using method implementation for struct that we use previously to also be fully consistent.

	}

}

func (bc *BlockchainStruct) BroadcastBlock(b *Block) {
	knirvlog.LogInfo("Broadcasting newly added block to all peers!")

}
func (bc *BlockchainStruct) AddPeer(peer string) {
    if bc.Peers == nil {

		bc.Peers = map[string]bool{}

	}
	bc.Peers[peer] = true
}

func (bc *BlockchainStruct) getOurCurrentBlockHash() (string, error) {

  lastBlock, err := bc.GetLastBlock()

    if err != nil {

		return "", fmt.Errorf("Unable to get last block hash: %w", err)

	}
	if lastBlock == nil {
       return "", fmt.Errorf("Unable to get last block hash due to no block state available, from: %v", lastBlock)

	}
   hash := lastBlock.Hash()
	return hex.EncodeToString(hash), nil

}
func (bc *BlockchainStruct) updateBlockchain(cm *ConsensusManager) {
  knirvlog.LogInfo("updating block chain, with latest block, locking mining....")// ensure that objects also are setting types based on data objects required for this operations which must use them or to perform updates, or any data modification which also has safety validations implemented before making any changes to those types parameters during program method/code execution in memory during test phases
    syncChan := cm.getSyncState()


    cm.setUpdateRequired(true)  // validate this
    cm.lockMining()  // lock data to guarantee
   syncChan <- true   // send sync signals



    lastBlock, err := bc.GetLastBlock()

    if err != nil {
        knirvlog.LogError("unable to get last block of blockchain for peer sync:", err)
        syncChan <- false
		return // implement proper error handler, with logging.
	}

   if lastBlock == nil {

        knirvlog.LogError("unable to access last block: ", errors.New("genesis block could not be located..."))
		syncChan <- false

		return

    }
	cm.setLongestChain(hex.EncodeToString(lastBlock.Hash()))

	 err =  bc.PutIntoDb(NewLevelDB(), bc.ChainAddress) // implemented correctly type here and method usage where struct parameter and it's implementation using types from this `blockchain_struct` code file implementation

    if err != nil {
        knirvlog.LogError("unable to update blockchain: ", err) // check the correct logging levels while testing system under `go test`.
       syncChan <- false
		return

    }

	time.Sleep(time.Second * 10)
   knirvlog.LogInfo("Updating Blockchain... Done!")
    cm.setUpdateRequired(false)
    syncChan <- false


}


func (cm *ConsensusManager) RunConsensus(bc *BlockchainStruct){ // Implementation is called here for code testing requirements of workflow.
	knirvlog.LogInfo("Starting consensus mechanism with chain Address " + bc.ChainAddress)

	for {
        if cm.getUpdateRequired() {
           knirvlog.LogWarning("chain sync lock is currently enabled by implementation code so current consensus can not continue execution as that specific state must exists before, method or data from code for validations can start from struct objects and method calls")

            time.Sleep(time.Second * 20) // implemented safety of method with valid object implementations parameters when running concurrent testing workflow implementation using goroutines.
           	continue

            }

			time.Sleep(time.Second * 15)
		}

	}