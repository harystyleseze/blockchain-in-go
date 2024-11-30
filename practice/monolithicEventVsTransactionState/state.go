package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// State represents the current state of the blockchain, including balances and the transaction pool.
type State struct {
	Balances  map[Account]uint
	txMempool []Tx
	dbFile    *os.File
}

// NewStateFromDisk creates a new state from the genesis file and transaction history.
func NewStateFromDisk() (*State, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Load genesis data
	genFilePath := filepath.Join(cwd, "../../database", "genesis.json")
	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, err
	}

	// Initialize balances from genesis data
	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	// Open transaction database (tx.db)
	txDbFilePath := filepath.Join(cwd, "../../database", "tx.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	// Initialize state
	state := &State{Balances: balances, txMempool: make([]Tx, 0), dbFile: f}

	// Replay all transactions from tx.db to reconstruct state
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Tx
		if err := json.Unmarshal(scanner.Bytes(), &tx); err != nil {
			return nil, err
		}

		// Apply the transaction to the state
		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}

// apply applies a transaction to the state, updating balances accordingly.
func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	// Validate that the sender has enough balance
	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance for %s", tx.From)
	}

	// Update balances
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}

// Add adds a transaction to the mempool and applies it to the state.
func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}
	s.txMempool = append(s.txMempool, tx)
	return nil
}

// Persist saves the transactions from the mempool to the tx.db file.
func (s *State) Persist() error {
	// Make a copy of the mempool because it will be modified
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	// Write each transaction to the file
	for i := 0; i < len(mempool); i++ {
		txJson, err := json.Marshal(mempool[i])
		if err != nil {
			return err
		}

		// Write transaction JSON with a newline
		if _, err = s.dbFile.Write(append(txJson, '\n')); err != nil {
			return err
		}

		// Remove the transaction from the mempool after writing
		s.txMempool = s.txMempool[1:]
	}

	return nil
}
