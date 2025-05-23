package blockcache

import (
	"github.com/ltcsuite/lnd/lntypes"
	"github.com/ltcsuite/lnd/multimutex"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/ltcutil"
	"github.com/ltcsuite/ltcd/wire"
	"github.com/ltcsuite/neutrino"
	"github.com/ltcsuite/neutrino/cache"
	"github.com/ltcsuite/neutrino/cache/lru"
)

// BlockCache is an lru cache for blocks.
type BlockCache struct {
	Cache     *lru.Cache[wire.InvVect, *neutrino.CacheableBlock]
	HashMutex *multimutex.Mutex[lntypes.Hash]
}

// NewBlockCache creates a new BlockCache with the given maximum capacity.
func NewBlockCache(capacity uint64) *BlockCache {
	return &BlockCache{
		Cache: lru.NewCache[wire.InvVect, *neutrino.CacheableBlock](
			capacity,
		),
		HashMutex: multimutex.NewMutex[lntypes.Hash](),
	}
}

// GetBlock first checks to see if the BlockCache already contains the block
// with the given hash. If it does then the block is fetched from the cache and
// returned. Otherwise the getBlockImpl function is used in order to fetch the
// new block and then it is stored in the block cache and returned.
func (bc *BlockCache) GetBlock(hash *chainhash.Hash,
	getBlockImpl func(hash *chainhash.Hash) (*wire.MsgBlock,
		error)) (*wire.MsgBlock, error) {

	bc.HashMutex.Lock(lntypes.Hash(*hash))
	defer bc.HashMutex.Unlock(lntypes.Hash(*hash))

	// Create an inv vector for getting the block.
	inv := wire.NewInvVect(wire.InvTypeWitnessBlock, hash)

	// Check if the block corresponding to the given hash is already
	// stored in the blockCache and return it if it is.
	cacheBlock, err := bc.Cache.Get(*inv)
	if err != nil && err != cache.ErrElementNotFound {
		return nil, err
	}
	if cacheBlock != nil {
		return cacheBlock.MsgBlock(), nil
	}

	// Fetch the block from the chain backends.
	block, err := getBlockImpl(hash)
	if err != nil {
		return nil, err
	}

	// Add the new block to blockCache. If the Cache is at its maximum
	// capacity then the LFU item will be evicted in favour of this new
	// block.
	_, err = bc.Cache.Put(
		*inv, &neutrino.CacheableBlock{
			Block: ltcutil.NewBlock(block),
		},
	)
	if err != nil {
		return nil, err
	}

	return block, nil
}
