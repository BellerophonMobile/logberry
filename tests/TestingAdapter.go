package tests

import (
	"bytes"
	"github.com/BellerophonMobile/logberry"
	"sync/atomic"
	"testing"
)

type TestingAdapter struct {
	std    *logberry.Root
	main   *logberry.Task
	t      *testing.T
	buffer bytes.Buffer
}

func (x *TestingAdapter) Write(p []byte) (int, error) {

	n, err := x.buffer.Write(p)

	if bytes.HasSuffix(p, []byte("\n")) {
		x.t.Logf("%v", x.buffer.String())
		x.buffer.Reset()
	}

	return n, err

}

var numtasks uint64

func newtaskuid() uint64 {
	// We have seen this atomic call cause problems on ARM...
	return atomic.AddUint64(&numtasks, 1) - 1
}

func SetStdTesting(t *testing.T) *TestingAdapter {

	adapter := &TestingAdapter{
		std:  logberry.Std,
		main: logberry.Main,
		t:    t,
	}

	logberry.Std = logberry.NewRoot(24)
	logberry.Std.AddOutputDriver(logberry.NewTextOutput(adapter, "testing"))
	logberry.Main = logberry.Std.Task("Test")

	return adapter

}

func (x *TestingAdapter) Stop() {
	logberry.Main.Success()
	logberry.Std.Stop()
	logberry.Main = x.main
	logberry.Std = x.std
}
