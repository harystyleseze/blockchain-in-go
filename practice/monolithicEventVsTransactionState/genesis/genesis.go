package genesis

import (
	"encoding/json"
	"os"
)

type Genesis struct {
	GenesisTime string          `json:"genesis_time"`
	ChainID     string          `json:"chain_id"`
	Balances    map[string]uint `json:"balances"`
}

// Load the genesis file and return a Genesis struct
func LoadGenesis(filename string) (*Genesis, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var genesis Genesis
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&genesis); err != nil {
		return nil, err
	}
	return &genesis, nil
}
