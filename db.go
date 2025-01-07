package main

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type DBClient struct {
	db     *sql.DB
	mu     sync.Mutex
	dbPath string
}

func NewDBClient(dbPath string) (*DBClient, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS blocks (
            hash TEXT PRIMARY KEY,
            prev_hash TEXT,
            timestamp INTEGER,
            data TEXT,
            nonce INTEGER
        );
    `)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &DBClient{db: db, dbPath: dbPath}, nil
}

func (d *DBClient) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.db.Close()
}
func (d *DBClient) InsertBlock(block *Block) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, err := d.db.Exec("INSERT INTO blocks (hash, prev_hash, timestamp, data, nonce) VALUES (?, ?, ?, ?, ?)",
		block.Hash, block.PrevHash, block.Timestamp, block.Data, block.Nonce)
	if err != nil {
		return fmt.Errorf("failed to insert block to db: %w", err)
	}
	return nil
}
func (d *DBClient) GetLastBlock() (*Block, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	var block Block
	row := d.db.QueryRow("SELECT hash, prev_hash, timestamp, data, nonce FROM blocks ORDER BY timestamp DESC LIMIT 1")
	err := row.Scan(&block.Hash, &block.PrevHash, &block.Timestamp, &block.Data, &block.Nonce)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil //return nil if no rows
		}
		return nil, fmt.Errorf("failed to get last block from db: %w", err)
	}
	return &block, nil
}
