package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Genesis struct {
	GenesisTime string         `json:"genesis_time"`
	ChainID     string         `json:"chain_id"`
	Balances    map[string]int `json:"balances"`
}

// Load the Genesis data from the file
func loadGenesis(filename string) (*Genesis, error) {
	// Check if the file exists
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, return nil so we can create new genesis data
			return nil, nil
		}
		return nil, err
	}

	// Read the file contents
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Check if the file is empty
	if len(data) == 0 {
		return nil, nil
	}

	// Parse the data into the Genesis struct
	var genesis Genesis
	err = json.Unmarshal(data, &genesis)
	if err != nil {
		return nil, err
	}

	return &genesis, nil
}

// Save the Genesis data to the file
func saveGenesis(filename string, genesis *Genesis) error {
	// Marshal the data into JSON with indentation for better readability
	data, err := json.MarshalIndent(genesis, "", "  ")
	if err != nil {
		return err
	}

	// Write the data to the file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Make a purchase if the user has enough balance
func makePurchase(genesis *Genesis, user string, amount int) bool {
	if balance, exists := genesis.Balances[user]; exists && balance >= amount {
		genesis.Balances[user] -= amount
		return true
	}
	return false
}

func main() {
	// Ensure the database directory exists, otherwise create it
	err := os.MkdirAll("../../database", os.ModePerm)
	if err != nil {
		log.Fatal("Error creating directory:", err)
	}

	// Load the genesis data from the file or create new data if file is empty or doesn't exist
	genesis, err := loadGenesis("../../database/genesis_file.json")
	if err != nil {
		log.Fatal("Error loading genesis data:", err)
	}

	// If genesis is nil, it means the file is either empty or doesn't exist, so create new data
	if genesis == nil {
		// Create new dummy data
		genesis = &Genesis{
			GenesisTime: time.Now().Format(time.RFC3339),
			ChainID:     "harystylesdb",
			Balances: map[string]int{
				"harystyles": 5000, // Add initial user "harystyles" with balance 5000
			},
		}

		// Save the new genesis data to the file
		err = saveGenesis("../../database/genesis_file.json", genesis)
		if err != nil {
			log.Fatal("Error saving genesis file:", err)
		}
		fmt.Println("Created new genesis data with dummy content.")
	} else {
		// If genesis file is not empty, add new users
		genesis.Balances["user1"] = 78888
		genesis.Balances["user2"] = 10000

		// Save the updated genesis data back to the file
		err = saveGenesis("../../database/genesis_file.json", genesis)
		if err != nil {
			log.Fatal("Error saving updated genesis file:", err)
		}
		fmt.Println("Added new users and balances to existing genesis data.")
	}

	// To interact with the genesis file, simulate a transaction by uncommenting line 121 to 132
	// // Simulate a purchase
	// itemCost := 1000
	// if makePurchase(genesis, "harystyles", itemCost) {
	// 	fmt.Println("Purchase successful!")
	// } else {
	// 	fmt.Println("Not enough tokens.")
	// }

	// // Print updated balances
	// fmt.Println("Updated Balances:")
	// for user, balance := range genesis.Balances {
	// 	fmt.Printf("%s: %d tokens\n", user, balance)
	// }
}
