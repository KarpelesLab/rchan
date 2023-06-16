[![GoDoc](https://godoc.org/github.com/KarpelesLab/rchan?status.svg)](https://godoc.org/github.com/KarpelesLab/rchan)

# rchan

ResponseChan, a simple method to get a uint64 value for a channel that can be used to receive a later response.

This can be used alongside time.Time to send requests to an async peer and connect responses back.

# usage

On one side:

```go
	id, ch := rchan.New()
	defer id.Release() // it is important to release the id after use
	sendRequest(id, ...) // will send the request via rpc/etc, response is expected to come back through a different channel

	select {
	case res <- ch:
		// got a response!
		return res, nil
	case <-ctx.Done():
		// expired
		return nil, ctx.Err()
	}
```

And another side receiving the response from the IPC/etc:

```go
	id := rchan.Id(response.id)
	id.SendTimeout(100*time.Millisecond, response.data) // this will safely send the response to the recipient, other methods are also available
```

Note that the sending side should always have a timeout in order to avoid
deadlocks from a race condition happening between the time when read times
out and the channel is released. Go's `chan` structure provides no good way
to deal with this, as for example closing the channel will cause writes to
panic without providing a way to check for exceptions.

This could be avoided by using a more complex structure but the goal of this
library is to provide the lightest possible way to use go channels. This race
condition case is unlikely enough so that a timeout will most likely work.
