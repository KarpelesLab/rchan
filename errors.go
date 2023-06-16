package rchan

import "errors"

var (
	ErrChanClosed = errors.New("This channel has already been closed")
)
