package main

import (
	"fmt"
	// Import the necessary packages
)

func main() {
	// File paths for genesis and transaction files
	// genesisFile := "../../database/genesis.json"
	// stateFile := "../../database/state.json"

	// Read the genesis file to get the initial state
	state, err := NewStateFromDisk()
	if err != nil {
		fmt.Printf("Error reading genesis file: %v\n", err)
		return
	}

	// Display initial state (from the genesis file)
	fmt.Println("Initial balances (from Genesis file):")
	for user, balance := range state.Balances {
		fmt.Printf("%s: %d tokens\n", user, balance)
	}

	// Perform a transaction: Transfer tokens between two users
	tx1 := Tx{
		From:  "andrej",
		To:    "babayaga",
		Value: 2000,
		Data:  "",
	}

	err = state.Add(tx1)
	if err != nil {
		fmt.Printf("Error during transaction: %v\n", err)
		return
	}

	// Persist the new state and transactions
	err = state.Persist()
	if err != nil {
		fmt.Printf("Error persisting state file: %v\n", err)
		return
	}

	// Display updated state after transaction
	fmt.Println("\nUpdated balances after transaction:")
	for user, balance := range state.Balances {
		fmt.Printf("%s: %d tokens\n", user, balance)
	}
}
