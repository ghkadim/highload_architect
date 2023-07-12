package closer

import (
	"sync/atomic"
)

type Closer interface {
	Close() error
}

func FromFunction(fn func() error) Closer {
	return &functionImpl{
		fn: fn,
	}
}

type functionImpl struct {
	fn     func() error
	closed atomic.Bool
}

func (f *functionImpl) Close() error {
	if f.closed.Swap(true) {
		return nil
	}
	return f.fn()
}

var (
	nop = &nopImpl{}
)

func Nop() Closer {
	return nop
}

type nopImpl struct{}

func (*nopImpl) Close() error {
	return nil
}
