package main

import (
	"KNIRVCHAIN-DEV/constants"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"time"
)

type Config struct {
	Port          uint64
	MinersAddress string
	DatabasePath  string
}

func init() {
	log.SetPrefix(constants.BLOCKCHAIN_NAME + ":")
}

func main() {
	var config Config
	chainCmdSet := flag.NewFlagSet("chain", flag.ExitOnError)
	chainPort := chainCmdSet.Uint64("port", 5000, "HTTP port to launch our blockchain server")
	chainMiner := chainCmdSet.String("miners_address", "", "Miners address to credit mining reward")
	dbPath := chainCmdSet.String("database_path", filepath.Join(".", "knirv.db"), "Filepath for saving chain's database")
	chainCmdSet.Parse(os.Args[1:])

	if chainCmdSet.Parsed() {
		if *chainMiner == "" || chainCmdSet.NFlag() == 0 {
			fmt.Println("Usage of chain subcommand: ")
			chainCmdSet.PrintDefaults()
			os.Exit(1)
		}
		config = Config{
			Port:          *chainPort,
			MinersAddress: *chainMiner,
			DatabasePath:  *dbPath, // reading dbPath flag argument set by command-line with type parameters correctly.
		}

		db := NewLevelDB() // Proper object instance based on implementation details from those structs.
		genesisBlock := NewBlock([]byte{}, 0, 0, nil)
		blockchain, err := NewBlockchain(genesisBlock, config.MinersAddress, db)
		if err != nil {
			log.Fatalf("Failed to create blockchain: %v", err)
		}
		defer db.Close()
		blockchain.Peers[blockchain.ChainAddress] = true
		bcs := NewBlockchainServer(config.Port, blockchain, config.MinersAddress) // Use configuration parameters passed by type instead of local hardcoded properties.
		go bcs.Start()
		go blockchain.ProofOfWorkMining(config.MinersAddress, blockchain.ConsensusManager)
		go blockchain.DialAndUpdatePeers()

		tp := NewTransactionPool()                // Create and initialize
		cm := NewConsensusManager(blockchain, tp) // Creates proper struct of a defined `ConsensusManager` that implements interface or method parameters of a certain struct and those required workflow steps that tests require using that new type variable that now is available by code design when implementations for type validation is being enforced with object implementation.

		go blockchain.updateBlockchain(cm)
		time.Sleep(20 * time.Second)
	}
}
