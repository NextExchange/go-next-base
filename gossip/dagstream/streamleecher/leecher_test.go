package streamleecher

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/NextSmartChain/go-next-base/gossip/dagstream"
	"github.com/NextSmartChain/go-next-base/hash"
	"github.com/NextSmartChain/go-next-base/inter/idx"
)

func TestLeecherNoDeadlocks(t *testing.T) {
	for try := 0; try < 10; try++ {
		testLeecherNoDeadlocks(t, 1+rand.Intn(500))
	}
}

type peerRequest struct {
	peer    string
	request dagstream.Request
}

func testLeecherNoDeadlocks(t *testing.T, maxPeers int) {
	requests := make(chan peerRequest, 1000)
	config := DefaultConfig()
	config.RecheckInterval = time.Millisecond * 5
	config.MinSessionRestart = 2 * time.Millisecond * 5
	config.MaxSessionRestart = 5 * time.Millisecond * 5
	config.BaseProgressWatchdog = 3 * time.Millisecond * 5
	config.Session.RecheckInterval = time.Millisecond
	epoch := idx.Epoch(1)
	leecher := New(epoch, rand.Intn(2) == 0, config, Callbacks{
		OnlyNotConnected: func(ids hash.Events) hash.Events {
			if rand.Intn(2) == 0 {
				return hash.Events{}
			}
			return ids
		},
		RequestChunk: func(peer string, r dagstream.Request) error {
			requests <- peerRequest{peer, r}
			return nil
		},
		Suspend: func(peer string) bool {
			return rand.Intn(10) == 0
		},
		PeerEpoch: func(peer string) idx.Epoch {
			return 1 + epoch/2 + idx.Epoch(rand.Intn(int(epoch*2)))
		},
	})
	terminated := false
	for i := 0; i < maxPeers*2; i++ {
		peer := strconv.Itoa(rand.Intn(maxPeers))
		coin := rand.Intn(100)
		if coin <= 50 {
			err := leecher.RegisterPeer(peer)
			if !terminated {
				require.NoError(t, err)
			}
		} else if coin <= 60 {
			err := leecher.UnregisterPeer(peer)
			if !terminated {
				require.NoError(t, err)
			}
		} else if coin <= 65 {
			epoch++
			leecher.OnNewEpoch(epoch)
		} else if coin <= 70 {
			leecher.ForceSyncing()
		} else {
			time.Sleep(time.Millisecond)
		}
		select {
		case req := <-requests:
			if rand.Intn(10) != 0 {
				err := leecher.NotifyChunkReceived(req.request.Session.ID, hash.FakeEvent(), rand.Intn(5) == 0)
				if !terminated {
					require.NoError(t, err)
				}
			}
		default:
		}
		if !terminated && rand.Intn(maxPeers*2) == 0 {
			terminated = true
			leecher.terminate()
		}
	}
	if !terminated {
		leecher.Stop()
	} else {
		leecher.wg.Wait()
	}
}
