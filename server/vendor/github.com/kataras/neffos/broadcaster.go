package neffos

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// async broadcaster, doesn't wait for a publish to complete to all clients before any
// next broadcast call.
type broadcaster struct {
	messages []Message
	mu       *sync.Mutex
	awaiter  unsafe.Pointer
}

func newBroadcaster() *broadcaster {
	ch := make(chan struct{})
	awaiter := unsafe.Pointer(&ch)
	return &broadcaster{
		mu:      new(sync.Mutex),
		awaiter: awaiter,
	}
}

func (b *broadcaster) getAwaiter() <-chan struct{} {
	ptr := atomic.LoadPointer(&b.awaiter)
	return *((*chan struct{})(ptr))
}

func (b *broadcaster) broadcast(msgs []Message) {
	b.mu.Lock()
	b.messages = msgs
	b.mu.Unlock()

	ch := make(chan struct{})
	old := atomic.SwapPointer(&b.awaiter, unsafe.Pointer(&ch))
	close(*(*chan struct{})(old))
}

func (b *broadcaster) waitUntilClosed(closeCh <-chan struct{}) (msgs []Message, ok bool) {
	ch := b.getAwaiter()
	b.mu.Unlock()
	select {
	case <-ch:
		msgs = b.messages[:]
		ok = true
	case <-closeCh:
	}

	b.mu.Lock()
	return
}
