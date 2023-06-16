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

And another side receiving the response from the IPC/Etc:

```go
	id := rchan.Id(response.id)
	id.Send(ctx, response.data) // this will safely send the response to the recipient, other methods are available
```
