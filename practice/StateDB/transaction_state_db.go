package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Define the structure of the blockchain state
type State struct {
	Balances map[string]int `json:"balances"`
}

// Read the genesis file and return the initial state
func readGenesisFile(filename string) (*State, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open genesis file: %v", err)
	}
	defer file.Close()

	var genesisState State
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&genesisState)
	if err != nil {
		return nil, fmt.Errorf("could not read genesis file: %v", err)
	}

	return &genesisState, nil
}

// Update the state and write it to the specified file
func updateState(filename string, state *State) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal state: %v", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write state to file: %v", err)
	}
	return nil
}

// Perform a transaction: transfer tokens from one user to another
func makeTransaction(state *State, from string, to string, amount int) error {
	// Check if the 'from' user exists
	if _, exists := state.Balances[from]; !exists {
		return fmt.Errorf("user %s does not exist", from)
	}

	// Ensure the 'from' user has enough tokens
	if state.Balances[from] < amount {
		return fmt.Errorf("insufficient balance for user %s", from)
	}

	// Perform the transaction
	state.Balances[from] -= amount
	state.Balances[to] += amount
	return nil
}

func main() {
	// File paths for genesis and state files
	genesisFile := "../../database/genesis_file.json"
	stateFile := "../../database/state.json"

	// Read the genesis file to get the initial state
	genesisState, err := readGenesisFile(genesisFile)
	if err != nil {
		fmt.Printf("Error reading genesis file: %v\n", err)
		return
	}

	// Create a copy of the genesis state to track the current state
	state := &State{
		Balances: genesisState.Balances,
	}

	// Display initial state (from the genesis file)
	fmt.Println("Initial balances (from Genesis file):")
	for user, balance := range state.Balances {
		fmt.Printf("%s: %d tokens\n", user, balance)
	}

	// Perform a transaction: Transfer tokens between two users
	err = makeTransaction(state, "user1", "harystyles", 200)
	if err != nil {
		fmt.Printf("Error during transaction: %v\n", err)
		return
	}

	// Update the state after the transaction and write to state.json
	err = updateState(stateFile, state)
	if err != nil {
		fmt.Printf("Error updating state file: %v\n", err)
		return
	}

	// Display updated state after transaction
	fmt.Println("\nUpdated balances after transaction:")
	for user, balance := range state.Balances {
		fmt.Printf("%s: %d tokens\n", user, balance)
	}

}
