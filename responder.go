package rchan

import (
	"sync"
	"sync/atomic"
)

var (
	rchanNextId uint64
	rchanMap    = make(map[Id]chan any)
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
