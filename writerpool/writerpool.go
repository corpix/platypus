package writerpool

import (
	"io"
	"sync"
)

// WriterPool represents a dynamic and iterable pool of Writer interfaces.
type WriterPool struct {
	mutex *sync.Mutex

	// FIXME: More effective solution is possible if we will use
	// some tree data-structure which is concurrency-friendly.
	writers []io.Writer
}

// Add adds a Writer to the pool.
func (wp *WriterPool) Add(c io.Writer) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	wp.writers = append(
		wp.writers,
		c,
	)
}

// Has returns true if Writer exists in the pool and false otherwise.
func (wp *WriterPool) Has(c io.Writer) bool {
	for _, v := range wp.writers {
		if v == c {
			return true
		}
	}
	return false
}

// Remove removes Writer from the pool if it exists
// and returns true in this case, otherwise it will be
// false.
func (wp *WriterPool) Remove(c io.Writer) bool {
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

// Iter performs a full-scan of the pool with fn
// invoked for every Writer in the pool.
func (wp *WriterPool) Iter(fn func(io.Writer)) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	for _, v := range wp.writers {
		fn(v)
	}
}

// New creates new WriterPool.
func New() *WriterPool {
	return &WriterPool{
		mutex:   &sync.Mutex{},
		writers: []io.Writer{},
	}
}
