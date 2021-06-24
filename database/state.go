package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// State used to control the state of the chain
type State struct {
	Balances        map[Account]uint
	txMempool       []Tx
	dbFile          *os.File
	latestBlockHash Hash
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

	f, err := os.OpenFile(filepath.Join(cwd, "database", "block.db"), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	state := &State{
		Balances:        balances.Balances,
		txMempool:       []Tx{},
		latestBlockHash: Hash{},
		dbFile:          f,
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var blockFs BlockFS
		blockFsJson := scanner.Bytes()
		err = json.Unmarshal(blockFsJson, &blockFs)

		if err != nil {
			return nil, err
		}

		if err := state.applyBlock(blockFs.Value); err != nil {
			return nil, err
		}
		state.latestBlockHash = blockFs.Key
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

// LatestBlockHash -> get the last snapshot of tx.db transactions
func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

// Close -> close file overwrite for State
func (s *State) Close() {
	s.dbFile.Close()
}

// applyBlock apply the transaction on blocks that exist
func (s *State) applyBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.apply(tx); err != nil {
			return err
		}
	}
	return nil
}

// AddBlock add transaction to the block
func (s *State) AddBlock(b Block) error {
	for _, tx := range b.TXs {
		if err := s.AddTx(tx); err != nil {
			return err
		}
	}
	return nil
}

// Add -> add new tx to mempool
func (s *State) AddTx(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

// Persist -> persist transactions into tx file
func (s *State) Persist() (Hash, error) {
	block := NewBlock(s.latestBlockHash, uint64(time.Now().Unix()), s.txMempool)
	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err 
	}
	blockFs := BlockFS{blockHash, block}
	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err 
	}
	fmt.Printf("Persisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	if _, err := s.dbFile.Write(append(blockFsJson, '\n')); err != nil {
		return Hash{}, err
	}

	s.latestBlockHash = blockHash
	s.txMempool = []Tx{}
	
	return s.latestBlockHash, nil 
}
