package database

type Account string

func NewAccount(value string) Account {
	return Account(value)
}

type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

// NewTx -> return a new Tx
func NewTx(from Account, to Account, value uint, data string) Tx {
	return Tx{from, to, value, data}
}

// IsReward -> check if Data field is eq to reward
func (t Tx) IsReward() bool {
	return t.Data == "reward"
}
