package tdag

import (
	"github.com/NextSmartChain/go-next-base/hash"
	"github.com/NextSmartChain/go-next-base/inter/dag"
)

type TestEvent struct {
	dag.MutableBaseEvent
	Name string
}

func (e *TestEvent) AddParent(id hash.Event) {
	parents := e.Parents()
	parents.Add(id)
	e.SetParents(parents)
}
