package database

import (
	"encoding/json"
	"io/ioutil"
)

type gen struct {
	Balances map[Account]uint `json:"balances"`
}

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
