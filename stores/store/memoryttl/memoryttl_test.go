package memoryttl

import (
	"testing"
	"time"

	"github.com/corpix/loggers/logger/logrus"
	logrusLogger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	jsonTime "github.com/cryptounicorns/platypus/time"
)

func TestMemoryTTL(t *testing.T) {
	samples := []struct {
		name   string
		input  map[string]interface{}
		output map[string]interface{}
		sleep  time.Duration
	}{
		{
			name:   "empty",
			input:  map[string]interface{}{},
			output: map[string]interface{}{},
		},
		{
			name:   "single",
			input:  map[string]interface{}{"foo": "bar"},
			output: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "multiple",
			input: map[string]interface{}{
				"foo": "bar",
				"bar": 1,
				"baz": nil,
			},
			output: map[string]interface{}{
				"foo": "bar",
				"bar": 1,
				"baz": nil,
			},
		},
		{
			name: "multiple sleep",
			input: map[string]interface{}{
				"foo": "bar",
				"bar": 1,
				"baz": nil,
			},
			output: map[string]interface{}{},
			sleep:  50 * time.Millisecond,
		},
	}

	var (
		log = logrus.New(logrusLogger.New())
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				var (
					store       *MemoryTTL
					initialized = map[string]interface{}{}
					filled      = map[string]interface{}{}
					err         error
				)

				store, err = New(
					Config{
						TTL:        jsonTime.Duration(10 * time.Millisecond),
						Resolution: jsonTime.Duration(5 * time.Millisecond),
					},
					log,
				)
				if err != nil {
					t.Error(err)
					return
				}
				defer store.Close()

				store.Iter(
					func(key string, value interface{}) {
						initialized[key] = value
					},
				)
				assert.EqualValues(
					t,
					map[string]interface{}{},
					initialized,
				)

				for k, v := range sample.input {
					store.Set(k, v)
				}

				time.Sleep(sample.sleep)

				store.Iter(
					func(key string, value interface{}) {
						filled[key] = value
					},
				)
				assert.EqualValues(
					t,
					sample.output,
					filled,
				)

				for k, v := range sample.output {
					vv, err := store.Get(k)
					if err != nil {
						t.Error(err)
						return
					}

					assert.EqualValues(
						t,
						v,
						vv,
					)
				}

				time.Sleep(1 * time.Second)
			},
		)
	}
}
