package rchan

import "errors"

// ErrChanClosed is returned by Send and SendTimeout when attempting to send
// to an Id that has already been released or was never valid.
var ErrChanClosed = errors.New("channel has already been closed")
