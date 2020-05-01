package neffos

import (
	"sync/atomic"
)

// waiterOnce is used on the server and client-side connections to describe the readiness of handling messages.
// For both sides if Reading is errored it returns the error back to the `waiterOnce#wait()`.
// For server-side:
// It waits until error from `OnConnected` (if exists) or first Write action (i.e `Connect` on `OnConnected`).
//
// For client-side:
// It waits until ACK is done, if server sent an error then it returns the error to the `Client#Dial`.
//
// See `Server#ServeHTTP`, `Conn#Connect`, `Conn#Write`, `Conn#sendClientACK` and `Conn#handleACK`.
type waiterOnce struct {
	locked *uint32
	ready  *uint32
	err    error
	// mu     sync.Mutex
	ch chan struct{}
}

func newWaiterOnce() *waiterOnce {
	return &waiterOnce{
		locked: new(uint32),
		ready:  new(uint32),
		ch:     make(chan struct{}),
	}
}

func (w *waiterOnce) isReady() bool {
	if w == nil {
		return true
	}

	return atomic.LoadUint32(w.ready) > 0
}

// waits and returns any error from the `unwait`,
// but if `unwait` called before `wait` then it returns immediately.
func (w *waiterOnce) wait() error {
	if w == nil {
		return nil
	}

	if w.isReady() {
		// println("waiter: wait() is Ready")
		return w.err // no need to wait.
	}

	if atomic.CompareAndSwapUint32(w.locked, 0, 1) {
		// println("waiter: lock")
		<-w.ch
	}

	return w.err
}

func (w *waiterOnce) unwait(err error) {
	if w == nil || w.isReady() {
		return
	}

	w.err = err
	// at any case mark it as ready for future `wait` call to exit immediately.
	atomic.StoreUint32(w.ready, 1)
	if atomic.CompareAndSwapUint32(w.locked, 1, 0) { // unlock once.
		close(w.ch)
	}
}
