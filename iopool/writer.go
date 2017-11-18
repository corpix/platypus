package iopool

import (
	"io"
	"sync"
)

// Writer represents a dynamic and iterable  of Writer interfaces.
type Writer struct {
	mutex *sync.Mutex

	// FIXME: More effective solution is possible if we will use
	// some tree data-structure which is concurrency-friendly.
	writers []io.Writer
}

// Add adds a Writer to the .
func (wp *Writer) Add(c io.Writer) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	wp.writers = append(
		wp.writers,
		c,
	)
}

// Has returns true if Writer exists in the  and false otherwise.
func (wp *Writer) Has(c io.Writer) bool {
	for _, v := range wp.writers {
		if v == c {
			return true
		}
	}
	return false
}

// Remove removes Writer from the  if it exists
// and returns true in this case, otherwise it will be
// false.
func (wp *Writer) Remove(c io.Writer) bool {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	for k, v := range wp.writers {
		if v == c {
			if k < len(wp.writers)-1 {
				wp.writers = append(
					wp.writers[0:k],
					wp.writers[k+1:]...,
				)
			} else {
				wp.writers = wp.writers[0:k]
			}
			return true
		}
	}

	return false
}

// Iter performs a full-scan of the  with fn
// invoked for every Writer in the .
func (wp *Writer) Iter(fn func(io.Writer)) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	for _, v := range wp.writers {
		fn(v)
	}
}

// NewWriter creates new Writer.
func NewWriter() *Writer {
	return &Writer{
		mutex:   &sync.Mutex{},
		writers: []io.Writer{},
	}
}
