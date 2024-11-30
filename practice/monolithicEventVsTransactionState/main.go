package main

import (
	"fmt"
	"log"

	// Import the necessary packages
	"blockchain-in-go/practice/monolithicEventVsTransactionState/state"
	"blockchain-in-go/practice/monolithicEventVsTransactionState/transaction"
)

func main() {
	// Read the genesis file and create the initial state
	state, err := state.NewStateFromDisk() // Use the NewStateFromDisk from the state package
	if err != nil {
		log.Fatalf("Error reading state from disk: %v", err)
	}

	// Display initial state (from the genesis file)
	fmt.Println("Initial balances (from Genesis file):")
	for user, balance := range state.Balances {
		fmt.Printf("%s: %d tokens\n", user, balance)
	}

	// Perform a transaction: Transfer tokens between two users
	tx1 := transaction.Tx{ // Referencing the Tx struct from the transaction package
		From:  "harystyles",
		To:    "okeke",
		Value: 150,
		Data:  "laptop fee",
	}

	// Add the transaction to the state
	err = state.Add(tx1)
	if err != nil {
		log.Printf("Error during transaction: %v\n", err)
		return
	}

	// Persist the new state and transactions
	err = state.Persist()
	if err != nil {
		log.Printf("Error persisting state file: %v\n", err)
		return
	}

	// Display updated state after transaction
	fmt.Println("\nUpdated balances after transaction:")
	for user, balance := range state.Balances {
		fmt.Printf("%s: %d tokens\n", user, balance)
	}
}
