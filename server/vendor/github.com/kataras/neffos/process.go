package neffos

import (
	"sync"
	"sync/atomic"
)

// processes is a collection of `process`.
type processes struct {
	entries map[string]*process
	locker  *sync.RWMutex
}

func newProcesses() *processes {
	return &processes{
		entries: make(map[string]*process),
		locker:  new(sync.RWMutex),
	}
}

func (p *processes) get(name string) *process {
	p.locker.RLock()
	entry := p.entries[name]
	p.locker.RUnlock()

	if entry == nil {
		entry = &process{
			v: new(uint32),
		}
		p.locker.Lock()
		p.entries[name] = entry
		p.locker.Unlock()
	}

	return entry
}

// process is used on connections on specific actions that needs to wait for an answer from the other side.
// Take for example the `Conn#handleMessage.tryNamespace` which waits for `Conn#askConnect` to finish on the specific namespace.
type process struct {
	v *uint32
}

// func (p *process) run() func() {
// 	p.start()
// 	return p.stop
// }

func (p *process) start() {
	for !atomic.CompareAndSwapUint32(p.v, 0, 1) {
		// if already started then wait to finish.
	}
}

func (p *process) stop() {
	atomic.StoreUint32(p.v, 0)
}

func (p *process) wait() {
	for p.isRunning() {
	}
}

// returns true if process didn't start yet or if stopped running.
func (p *process) isRunning() bool {
	return atomic.LoadUint32(p.v) > 0
}
