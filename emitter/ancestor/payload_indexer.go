package ancestor

import (
	"github.com/NextSmartChain/go-next-base/hash"
	"github.com/NextSmartChain/go-next-base/inter/dag"
	"github.com/NextSmartChain/go-next-base/utils/wlru"
)

type PayloadIndexer struct {
	payloadLamports *wlru.Cache
}

func NewPayloadIndexer(cacheSize int) *PayloadIndexer {
	cache, _ := wlru.New(uint(cacheSize), cacheSize)
	return &PayloadIndexer{cache}
}

func (h *PayloadIndexer) ProcessEvent(event dag.Event, payloadMetric Metric) {
	maxParentsPayloadMetric := Metric(0)
	for _, p := range event.Parents() {
		parentMetric := h.GetMetricOf(p)
		if maxParentsPayloadMetric < parentMetric {
			maxParentsPayloadMetric = parentMetric
		}
	}
	if maxParentsPayloadMetric != 0 || payloadMetric != 0 {
		h.payloadLamports.Add(event.ID(), maxParentsPayloadMetric+payloadMetric, 1)
	}
}

func (h *PayloadIndexer) GetMetricOf(id hash.Event) Metric {
	parentMetric, ok := h.payloadLamports.Get(id)
	if !ok {
		return 0
	}
	return parentMetric.(Metric)
}

func (h *PayloadIndexer) SearchStrategy() SearchStrategy {
	return NewMetricStrategy(h.GetMetricOf)
}
