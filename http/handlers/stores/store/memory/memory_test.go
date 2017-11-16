package memory

import (
	"testing"

	"github.com/corpix/loggers/logger/logrus"
	logrusLogger "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	samples := []struct {
		name string
		data map[string]interface{}
	}{
		{
			name: "empty",
			data: map[string]interface{}{},
		},
		{
			name: "single",
			data: map[string]interface{}{"foo": "bar"},
		},
		{
			name: "multiple",
			data: map[string]interface{}{
				"foo": "bar",
				"bar": 1,
				"baz": nil,
			},
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
					store       *Memory
					initialized = map[string]interface{}{}
					filled      = map[string]interface{}{}
					err         error
				)

				store, err = New(Config{}, log)
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

				for k, v := range sample.data {
					store.Set(k, v)
				}

				store.Iter(
					func(key string, value interface{}) {
						filled[key] = value
					},
				)
				assert.EqualValues(
					t,
					sample.data,
					filled,
				)

				for k, v := range sample.data {
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
			},
		)
	}
}
