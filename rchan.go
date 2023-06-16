package rchan

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

var (
	rchanNextId uint64
	rchanMap    = make(map[Id]chan<- any)
	rchanMapLk  sync.RWMutex
)

type Id uint64

// New creates a new rchan and returns the associated chan to read from to obtain the delayed response.
func New() (Id, <-chan any) {
	c := make(chan any)

	rchanMapLk.Lock()
	defer rchanMapLk.Unlock()

	for {
		id := Id(atomic.AddUint64(&rchanNextId, 1))
		if id == 0 {
			// we never return id==0 so the using lib can use id=0 to mean "no response expected"
			continue
		}
		if _, f := rchanMap[id]; f {
			// this will probably never happen unless Release() isn't called timely
			continue
		}
		rchanMap[id] = c
		return id, c
	}
}

// Release will release the given object from memory, future calls to C() on this object
// will return nil.
func (r Id) Release() {
	rchanMapLk.Lock()
	defer rchanMapLk.Unlock()

	delete(rchanMap, r)
}

// C will return the channel associated with a rchan object, or nil, and can be used
// to send the response to the initiator
func (r Id) C() chan<- any {
	if r == 0 {
		return nil
	}

	rchanMapLk.RLock()
	defer rchanMapLk.RUnlock()

	if c, ok := rchanMap[r]; ok {
		return c
	}
	return nil
}

// Send will send the given value in the given channel unless the context expires first
func (r *Id) Send(ctx context.Context, v any) error {
	c := r.C()
	if c == nil {
		return ErrChanClosed
	}

	select {
	case c <- v:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// SendTimeout will perform a send that will expire after a given time, using a timer
// rather than contexts in order to limit weight
func (r *Id) SendTimeout(max time.Duration, v any) error {
	c := r.C()
	if c == nil {
		return ErrChanClosed
	}

	t := time.NewTimer(max)
	defer t.Stop()

	select {
	case c <- v:
		return nil
	case <-t.C:
		// we use the error from context despite not using a context for consistency, and definitely not because I was too lazy to add another error to errors.go
		return context.DeadlineExceeded
	}
}
