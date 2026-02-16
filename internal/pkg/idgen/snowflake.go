package idgen

import (
	"sync"
	"time"
)

const (
	workerBits        = int64(10)
	sequenceBits      = int64(12)
	maxWorkerID       = int64(-1) ^ (int64(-1) << workerBits)
	sequenceMask      = int64(-1) ^ (int64(-1) << sequenceBits)
	workerShift       = sequenceBits
	timestampShift    = sequenceBits + workerBits
	defaultWorkerID   = int64(1)             // 默认worker ID
	defaultEpochMilli = int64(1767225600000) // 2026-01-01 00:00:00 UTC
)

type generator struct {
	mu       sync.Mutex
	workerID int64
	epochMs  int64
	lastMs   int64
	seq      int64
}

var (
	globalMu sync.RWMutex
	global   *generator
)

func Init() error {
	return New(defaultWorkerID, defaultEpochMilli)
}

func New(workerID, epochMs int64) error {
	if workerID < 0 || workerID > maxWorkerID {
		return ErrInvalidWorkerID
	}
	if epochMs <= 0 {
		epochMs = defaultEpochMilli
	}

	g := &generator{
		workerID: workerID,
		epochMs:  epochMs,
	}

	globalMu.Lock()
	global = g
	globalMu.Unlock()
	return nil
}

func NextID() int64 {
	g := getOrInitDefault()
	return g.nextID()
}

func getOrInitDefault() *generator {
	globalMu.RLock()
	g := global
	globalMu.RUnlock()
	if g != nil {
		return g
	}

	globalMu.Lock()
	defer globalMu.Unlock()
	if global == nil {
		global = &generator{
			workerID: defaultWorkerID,
			epochMs:  defaultEpochMilli,
		}
	}
	return global
}

func (g *generator) nextID() int64 {
	g.mu.Lock()
	defer g.mu.Unlock()

	nowMs := nowMilli()
	if nowMs < g.lastMs {
		// Keep monotonic output when wall clock moves backward.
		nowMs = g.lastMs
	}

	if nowMs == g.lastMs {
		g.seq = (g.seq + 1) & sequenceMask
		if g.seq == 0 {
			nowMs = g.waitNextMillis(nowMs)
		}
	} else {
		g.seq = 0
	}

	g.lastMs = nowMs
	if nowMs < g.epochMs {
		nowMs = g.epochMs
	}

	return ((nowMs - g.epochMs) << timestampShift) | (g.workerID << workerShift) | g.seq
}

func (g *generator) waitNextMillis(lastMs int64) int64 {
	nowMs := nowMilli()
	for nowMs <= lastMs {
		time.Sleep(time.Millisecond)
		nowMs = nowMilli()
	}
	return nowMs
}

func nowMilli() int64 {
	return time.Now().UnixMilli()
}
