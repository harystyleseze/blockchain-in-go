package main

// Transaction represents a transfer of value from one account to another.
type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"` // "reward" or other data
}

// IsReward checks if the transaction is a reward type.
func (t Tx) IsReward() bool {
	return t.Data == "reward"
}
