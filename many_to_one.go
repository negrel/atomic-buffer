package ringo

import (
	"sync/atomic"
)

type manyToOne buffer

// ManyToOne return an efficient buffer with the given capacity.
// The buffer is safe for one reader and multiple writer.
func ManyToOne(capacity uint32) Buffer {
	return &oneToOne{
		head:     ^uint64(0),
		buffer:   make([]Generic, capacity),
		capacity: uint64(capacity),
	}
}

func (mto *manyToOne) Cap() uint32 {
	return uint32(mto.capacity)
}

// Push the given data to the buffer.
func (mto *manyToOne) Push(data Generic) {
	head := atomic.AddUint64(&mto.head, 1)
	index := head % mto.capacity

	box := box{
		index: head,
		data:  data,
	}

	atomic.SwapPointer(&mto.buffer[index], Generic(&box))
}

// Push the given data to the buffer and return if
// the data is valid.
func (mto *manyToOne) Shift() (Generic, bool) {
	i := mto.tail % mto.capacity
	mto.tail++

	box := (*box)(mto.buffer[i])

	if box == nil {
		return nil, false
	}

	return box.data, box.index > mto.tail
}
