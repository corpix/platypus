package writerpool

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testConn uint8

func (t *testConn) Read([]byte) (int, error)         { panic("not implemented") }
func (t *testConn) Write([]byte) (int, error)        { panic("not implemented") }
func (t *testConn) Close() error                     { panic("not implemented") }
func (t *testConn) LocalAddr() net.Addr              { panic("not implemented") }
func (t *testConn) RemoteAddr() net.Addr             { panic("not implemented") }
func (t *testConn) SetDeadline(time.Time) error      { panic("not implemented") }
func (t *testConn) SetReadDeadline(time.Time) error  { panic("not implemented") }
func (t *testConn) SetWriteDeadline(time.Time) error { panic("not implemented") }

func newTestConn() *testConn { return new(testConn) }

func TestWriterPoolRemove(t *testing.T) {
	var (
		connections = []io.Writer{
			newTestConn(),
			newTestConn(),
			newTestConn(),
			newTestConn(),
		}
	)

	samples := []struct {
		name        string
		connections []io.Writer
		removed     bool
		remove      []io.Writer
		result      []io.Writer
	}{
		{
			name: "OneToOne",
			connections: []io.Writer{
				connections[0],
			},
			removed: true,
			remove:  []io.Writer{connections[0]},
			result:  []io.Writer{},
		},
		{
			name: "OneToMultiple",
			connections: []io.Writer{
				connections[0],
				connections[1],
				connections[2],
			},
			removed: true,
			remove:  []io.Writer{connections[0]},
			result: []io.Writer{
				connections[1],
				connections[2],
			},
		},
		{
			name: "InTheMiddle",
			connections: []io.Writer{
				connections[0],
				connections[1],
				connections[2],
			},
			removed: true,
			remove:  []io.Writer{connections[1]},
			result: []io.Writer{
				connections[0],
				connections[2],
			},
		},
		{
			name: "LastInTheTail",
			connections: []io.Writer{
				connections[0],
				connections[1],
				connections[2],
			},
			removed: true,
			remove:  []io.Writer{connections[2]},
			result: []io.Writer{
				connections[0],
				connections[1],
			},
		},
		{
			name: "NotFound",
			connections: []io.Writer{
				connections[0],
				connections[1],
				connections[2],
			},
			removed: false,
			remove:  []io.Writer{newTestConn()},
			result: []io.Writer{
				connections[0],
				connections[1],
				connections[2],
			},
		},
	}

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				pool := New()
				for _, v := range sample.connections {
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
