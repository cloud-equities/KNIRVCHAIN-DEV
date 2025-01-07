package main

import (
	"KNIRVCHAIN-DEV/constants"
	knirvlog "KNIRVCHAIN-DEV/log"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func init() {
	log.SetPrefix(constants.BLOCKCHAIN_NAME + ":")
}

func main() {
	chainCmdSet := flag.NewFlagSet("chain", flag.ExitOnError)
	chainPort := chainCmdSet.Uint64("port", 5000, "HTTP port to launch our blockchain server")
	chainMiner := chainCmdSet.String("miners_address", "", "Miners address to credit mining reward")
	//remoteNode := chainCmdSet.String("remote_node", "", "Remote Node from where the blockchain will be synced")
	dbPath := chainCmdSet.String("database_path", filepath.Join(".", "knirv.db"), "Filepath for saving chain's database") // Filepath logic in command args
	chainCmdSet.Parse(os.Args[1:])                                                                                        // Use cli data validation

	if chainCmdSet.Parsed() {

		if *chainMiner == "" || chainCmdSet.NFlag() == 0 {
			fmt.Println("Usage of chain subcommand: ")
			chainCmdSet.PrintDefaults() // output format info if arguments are wrong.
			os.Exit(1)                  // end workflow logic since the wrong data was passed into process
		}

		db, err := NewDBClient(*dbPath) // load a persistent object implementation. and verify it returns correctly before proceeding using error testing
		if err != nil {
			knirvlog.FatalError("Unable to create database client: ", err) // use testing scope, from project logger to indicate failed init, where required parameters, to methods or object, data declarations with correct error message formatting type declarations.

			os.Exit(1)
		}

		defer db.Close()

		chainAddress := "http://127.0.0.1:" + strconv.Itoa(int(*chainPort))
		genesisBlock := Block{}
		blockchain1 := NewBlockchain(&genesisBlock, chainAddress, db) // using chain parameter to verify what will happen with local instances and that types do not affect the code from crashing when using those type values as arguments.

		blockchain1.Peers[blockchain1.ChainAddress] = true

		bcs := NewBlockchainServer(*chainPort, blockchain1, chainAddress)
		go bcs.Start() // Test HTTP interface components.

		go bcs.BlockchainPtr.ProofOfWorkMining(*chainMiner) // all blockchain data components access logic types.

		go bcs.BlockchainPtr.DialAndUpdatePeers() // use that value that we passed here for methods logic in their execution path.
		go bcs.BlockchainPtr.RunConsensus()       // also make sure your core logic implementation test logic paths also work here when used with data being implemented via parameters and their method calls
		time.Sleep(20 * time.Second)              // halt system before exit. and if required log this.

	}

}
