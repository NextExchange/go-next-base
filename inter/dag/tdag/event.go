package tdag

import (
	"github.com/NextExchange/go-next-base/hash"
	"github.com/NextExchange/go-next-base/inter/dag"
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
