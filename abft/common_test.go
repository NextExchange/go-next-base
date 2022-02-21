package abft

import (
	"github.com/NextSmartChain/go-next-base/inter/idx"
	"github.com/NextSmartChain/go-next-base/inter/pos"
	"github.com/NextSmartChain/go-next-base/kvdb"
	"github.com/NextSmartChain/go-next-base/kvdb/memorydb"
	"github.com/NextSmartChain/go-next-base/orion"
	"github.com/NextSmartChain/go-next-base/utils/adapters"
	"github.com/NextSmartChain/go-next-base/vecfc"
)

type applyBlockFn func(block *orion.Block) *pos.Validators

// TestOrion extends Orion for tests.
type TestOrion struct {
	*IndexedOrion

	blocks map[idx.Block]*orion.Block

	applyBlock applyBlockFn
}

// FakeOrion creates empty abft with mem store and equal weights of nodes in genesis.
func FakeOrion(nodes []idx.ValidatorID, weights []pos.Weight, mods ...memorydb.Mod) (*TestOrion, *Store, *EventStore) {
	validators := make(pos.ValidatorsBuilder, len(nodes))
	for i, v := range nodes {
		if weights == nil {
			validators[v] = 1
		} else {
			validators[v] = weights[i]
		}
	}

	openEDB := func(epoch idx.Epoch) kvdb.DropableStore {
		return memorydb.New()
	}
	crit := func(err error) {
		panic(err)
	}
	store := NewStore(memorydb.New(), openEDB, crit, LiteStoreConfig())

	err := store.ApplyGenesis(&Genesis{
		Validators: validators.Build(),
		Epoch:      FirstEpoch,
	})
	if err != nil {
		panic(err)
	}

	input := NewEventStore()

	config := LiteConfig()
	lch := NewIndexedOrion(store, input, &adapters.VectorToDagIndexer{vecfc.NewIndex(crit, vecfc.LiteConfig())}, crit, config)

	extended := &TestOrion{
		IndexedOrion: lch,
		blocks:          map[idx.Block]*orion.Block{},
	}

	blockIdx := idx.Block(0)

	err = extended.Bootstrap(orion.ConsensusCallbacks{
		BeginBlock: func(block *orion.Block) orion.BlockCallbacks {
			blockIdx++
			return orion.BlockCallbacks{
				EndBlock: func() (sealEpoch *pos.Validators) {
					// track blocks
					extended.blocks[blockIdx] = block
					if extended.applyBlock != nil {
						return extended.applyBlock(block)
					}
					return nil
				},
			}
		},
	})
	if err != nil {
		panic(err)
	}

	return extended, store, input
}
