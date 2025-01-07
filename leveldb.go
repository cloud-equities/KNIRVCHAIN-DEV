package main

import (
	"encoding/json"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type LevelDB struct {
	Client *leveldb.DB
}

func NewDBClient(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, &opt.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to open leveldb: %w", err)
	}
	return &LevelDB{Client: db}, nil
}

func (db *LevelDB) SaveLastBlock(block interface{}, address string) error {
	blockBytes, err := json.Marshal(block)
	if err != nil {
		return err
	}
	err = db.Client.Put([]byte(address+"last_block"), blockBytes, &opt.WriteOptions{})
	if err != nil {
		return fmt.Errorf("failed to put data into database with key last_block: %w", err)
	}
	return nil
}

func (db *LevelDB) LoadLastBlock(address string) (interface{}, error) {
	data, err := db.Client.Get([]byte(address+"last_block"), &opt.ReadOptions{})
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, fmt.Errorf("no block found: %w", err)
		}
		return nil, fmt.Errorf("failed to get data from database with key last_block: %w", err)
	}

	block := new(Block)
	err = json.Unmarshal(data, block)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return block, nil
}

func (db *LevelDB) KeyExists(address string) (bool, error) {
	_, err := db.Client.Get([]byte(address), &opt.ReadOptions{})
	if err != nil {
		if err == leveldb.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to get data from database with key blockchain: %w", err)

	}
	return true, nil
}

func (db *LevelDB) GetBlockchain(address string) (interface{}, error) {
	data, err := db.Client.Get([]byte(address), &opt.ReadOptions{})
	if err != nil {
		if err == leveldb.ErrNotFound {
			return nil, fmt.Errorf("no blockchain found: %w", err)
		}
		return nil, fmt.Errorf("failed to get data from database with key blockchain: %w", err)
	}

	blockchain := new(BlockchainStruct)
	err = json.Unmarshal(data, blockchain)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return blockchain, nil
}

// new and correct implementations using object state to implement workflow, for types. where those data types and methods from structs/object should use those other structs and object methods from same program.

func (db *LevelDB) PutIntoDb(blockchain interface{}, address string) error { // this method receives implementations to work with data types to do saving using interfaces on workflow steps implementation, during runtime where data persists for software validation for the required logic and parameters for the tests workflow with those project related requirements when testing an application under methods executions when methods calls also use data from structs to save in `leveldb`.
	data, err := json.Marshal(blockchain) // validates if those types can marshal during method executions before it can save the new state object information into persistence leveldb database when calling during that workflow as a local workflow of object states and their properties
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = db.Client.Put([]byte(address), data, &opt.WriteOptions{}) // actually saves the new properties data/or values or object implementations to the databse, once validation pass without problems.

	if err != nil {
		return fmt.Errorf("failed to put data into database with key blockchain: %w", err)

	}
	return nil
}
func (db *LevelDB) Close() error {

	return db.Client.Close()

}
