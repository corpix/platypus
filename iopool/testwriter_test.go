package iopool

type testWriter uint8

func (t *testWriter) Write([]byte) (int, error) { panic("not implemented") }
func (t *testWriter) Close() error              { panic("not implemented") }

func newTestWriter() *testWriter { return new(testWriter) }
