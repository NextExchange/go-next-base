package abft

import (
	"github.com/NextSmartChain/go-next-base/hash"
	"github.com/NextSmartChain/go-next-base/inter/dag"
)

// EventSource is a callback for getting events from an external storage.
type EventSource interface {
	HasEvent(hash.Event) bool
	GetEvent(hash.Event) dag.Event
}
