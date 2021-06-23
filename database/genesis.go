package database

import (
	"encoding/json"
	"io/ioutil"
)

// gen -> load all the content from the genesis file
type gen struct {
	Balances map[Account]uint `json:"balances"`
}

// loadgen -> load genesis file content
func loadgen(path string) (gen, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return gen{}, err
	}

	var loadGen gen
	err = json.Unmarshal(content, &loadGen)
	if err != nil {
		return gen{}, err
	}

	return loadGen, nil
}
