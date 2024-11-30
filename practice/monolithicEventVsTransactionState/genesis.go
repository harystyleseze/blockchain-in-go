package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Genesis holds the initial blockchain state and metadata.
type Genesis struct {
	GenesisTime string           `json:"genesis_time"`
	ChainID     string           `json:"chain_id"`
	Balances    map[Account]uint `json:"balances"`
}

// loadGenesis loads the genesis data from a file.
func loadGenesis(filename string) (*Genesis, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open genesis file: %v", err)
	}
	defer file.Close()

	var gen Genesis
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&gen); err != nil {
		return nil, fmt.Errorf("could not read genesis file: %v", err)
	}
	return &gen, nil
}
