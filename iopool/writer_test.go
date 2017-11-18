package iopool

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWriterRemove(t *testing.T) {
	var (
		writers = []io.Writer{
			newTestWriter(),
			newTestWriter(),
			newTestWriter(),
			newTestWriter(),
		}
	)

	samples := []struct {
		name    string
		writers []io.Writer
		removed bool
		remove  []io.Writer
		result  []io.Writer
	}{
		{
			name: "OneToOne",
			writers: []io.Writer{
				writers[0],
			},
			removed: true,
			remove:  []io.Writer{writers[0]},
			result:  []io.Writer{},
		},
		{
			name: "OneToMultiple",
			writers: []io.Writer{
				writers[0],
				writers[1],
				writers[2],
			},
			removed: true,
			remove:  []io.Writer{writers[0]},
			result: []io.Writer{
				writers[1],
				writers[2],
			},
		},
		{
			name: "InTheMiddle",
			writers: []io.Writer{
				writers[0],
				writers[1],
				writers[2],
			},
			removed: true,
			remove:  []io.Writer{writers[1]},
			result: []io.Writer{
				writers[0],
				writers[2],
			},
		},
		{
			name: "LastInTheTail",
			writers: []io.Writer{
				writers[0],
				writers[1],
				writers[2],
			},
			removed: true,
			remove:  []io.Writer{writers[2]},
			result: []io.Writer{
				writers[0],
				writers[1],
			},
		},
		{
			name: "NotFound",
			writers: []io.Writer{
				writers[0],
				writers[1],
				writers[2],
			},
			removed: false,
			remove:  []io.Writer{newTestWriter()},
			result: []io.Writer{
				writers[0],
				writers[1],
				writers[2],
			},
		},
	}

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				pool := NewWriter()
				for _, v := range sample.writers {
					pool.Add(v)
				}

				for _, v := range sample.remove {
					assert.Equal(
						t,
						sample.removed,
						pool.Remove(v),
						"remove should be true if item is removed",
					)
				}

				result := []io.Writer{}
				pool.Iter(
					func(c io.Writer) {
						result = append(result, c)
					},
				)

				assert.Equal(
					t,
					sample.result,
					result,
					"should be same as result after items removed",
				)
			},
		)
	}
}
