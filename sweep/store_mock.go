package sweep

import (
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/wire"
)

// MockSweeperStore is a mock implementation of sweeper store. This type is
// exported, because it is currently used in nursery tests too.
type MockSweeperStore struct {
	ourTxes map[chainhash.Hash]struct{}
}

// NewMockSweeperStore returns a new instance.
func NewMockSweeperStore() *MockSweeperStore {
	return &MockSweeperStore{
		ourTxes: make(map[chainhash.Hash]struct{}),
	}
}

// IsOurTx determines whether a tx is published by us, based on its
// hash.
func (s *MockSweeperStore) IsOurTx(hash chainhash.Hash) (bool, error) {
	_, ok := s.ourTxes[hash]
	return ok, nil
}

// NotifyPublishTx signals that we are about to publish a tx.
func (s *MockSweeperStore) NotifyPublishTx(tx *wire.MsgTx) error {
	txHash := tx.TxHash()
	s.ourTxes[txHash] = struct{}{}

	return nil
}

// ListSweeps lists all the sweeps we have successfully published.
func (s *MockSweeperStore) ListSweeps() ([]chainhash.Hash, error) {
	var txns []chainhash.Hash
	for tx := range s.ourTxes {
		txns = append(txns, tx)
	}

	return txns, nil
}

// Compile-time constraint to ensure MockSweeperStore implements SweeperStore.
var _ SweeperStore = (*MockSweeperStore)(nil)
