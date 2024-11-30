package transaction

type Tx struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value uint   `json:"value"`
	Data  string `json:"data"`
}

// Check if the transaction is a reward
func (t Tx) IsReward() bool {
	return t.Data == "reward"
}
