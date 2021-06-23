package database

import (
	"bufio"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Snapshot [32]byte

type State struct {
	Balances  map[Account]uint
	txMempool []Tx
	snapshot  Snapshot
	dbFile    *os.File
	// latestBlockHash	Hash
}

//TODO: refactor to Hash instead of snapshot
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
		snapshot:  Snapshot{},
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

	err = state.doSnapshot()
	if err != nil {
		return nil, err
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

// LatestSnapshot -> get the last snapshot of tx.db transactions
func (s *State) LatestSnapshot() Snapshot {
	return s.snapshot
}

// doSnapshot -> create sha256 snapshot with tx.db content
func (s *State) doSnapshot() error {
	_, err := s.dbFile.Seek(0, 0)
	if err != nil {
		return err
	}

	txsData, err := ioutil.ReadAll(s.dbFile)
	if err != nil {
		return err
	}

	s.snapshot = sha256.Sum256(txsData)
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
func (s *State) Persist() (Snapshot, error) {
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {
		txJson, err := json.Marshal(mempool[i])
		if err != nil {
			return Snapshot{}, err
		}

		fmt.Println("persisting")
		fmt.Printf("\t%s\n", txJson)
		if _, err := s.dbFile.Write(append(txJson, '\n')); err != nil {
			return Snapshot{}, err
		}

		err = s.doSnapshot()
		if err != nil {
			return Snapshot{}, err
		}
		fmt.Printf("new db snapshot: %x\n", s.snapshot)

		s.txMempool = s.txMempool[1:]
	}

	return s.snapshot, nil
}
