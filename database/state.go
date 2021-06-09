package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Tx
	dbFile    *os.File
}

// NewStateFromDisk -> load all state information
func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	gen, err := loadgen(filepath.Join(cwd, "database", "gen.json"))
	if err != nil {
		return nil, err
	}

	var balances State
	balances.Balances = gen.Balances

	f, err := os.OpenFile(filepath.Join(cwd, "database", "tx.db"), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	state := &State{
		Balances:  balances.Balances,
		txMempool: []Tx{},
		dbFile:    f,
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)
		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}
	return state, nil
}

// apply -> apply tx to a certain State
func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From] < tx.Value {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

// Close -> close file overwrite for State
func (s *State) Close() {
	s.dbFile.Close()
}

// Add -> add new tx to mempool
func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

// Persist -> persist transactions into tx file
func (s *State) Persist() error {
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {
		txJson, err := json.Marshal(mempool[i])
		if err != nil {
			return err
		}

		if _, err := s.dbFile.Write(append(txJson, '\n')); err != nil {
			return err
		}

		s.txMempool = s.txMempool[1:]
	}

	return nil
}
