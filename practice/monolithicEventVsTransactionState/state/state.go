package state

import (
	"blockchain-in-go/practice/monolithicEventVsTransactionState/genesis"
	"blockchain-in-go/practice/monolithicEventVsTransactionState/transaction"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type State struct {
	Balances  map[string]uint
	TxMempool []transaction.Tx
	DbFile    *os.File
}

// Function to read the genesis file and create the initial state
func NewStateFromDisk() (*State, error) {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Load genesis file
	genFilePath := filepath.Join(cwd, "database", "genesis.json")
	gen, err := genesis.LoadGenesis(genFilePath)
	if err != nil {
		return nil, err
	}

	// Set up initial balances based on the genesis file
	balances := make(map[string]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	// Open transaction DB file
	txDbFilePath := filepath.Join(cwd, "database", "tx.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	// Initialize state
	state := &State{
		Balances:  balances,
		TxMempool: make([]transaction.Tx, 0),
		DbFile:    f,
	}

	// Rebuild the state by reading transactions from the transaction DB
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var tx transaction.Tx
		if err := json.Unmarshal(scanner.Bytes(), &tx); err != nil {
			return nil, err
		}

		// Apply the transaction to the state
		if err := state.Apply(tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}

// Apply a transaction to update the state
func (s *State) Apply(tx transaction.Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance for %s", tx.From)
	}

	// Apply regular transaction
	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value
	return nil
}

// Add a new transaction to the state
func (s *State) Add(tx transaction.Tx) error {
	if err := s.Apply(tx); err != nil {
		return err
	}
	s.TxMempool = append(s.TxMempool, tx)
	return nil
}

// Persist the transactions to disk
func (s *State) Persist() error {
	mempool := make([]transaction.Tx, len(s.TxMempool))
	copy(mempool, s.TxMempool)

	for _, tx := range mempool {
		txJson, err := json.Marshal(tx)
		if err != nil {
			return err
		}

		if _, err := s.DbFile.Write(append(txJson, '\n')); err != nil {
			return err
		}

		// Remove the transaction from the mempool after persisting
		s.TxMempool = s.TxMempool[1:]
	}
	return nil
}
