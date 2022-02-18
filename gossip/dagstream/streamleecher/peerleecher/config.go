package peerleecher

import (
	"time"

	"github.com/NextSmartChain/go-next-base/inter/dag"
)

type EpochDownloaderConfig struct {
	RecheckInterval        time.Duration
	DefaultChunkSize       dag.Metric
	ParallelChunksDownload int
}
