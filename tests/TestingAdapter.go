package tests

import (
	"github.com/BellerophonMobile/logberry"
	"bytes"
	"testing"
	"sync/atomic"
)

type TestingAdapter struct {
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

func SetStdTesting(t *testing.T) {

	logberry.Std = logberry.NewRoot(24)
	logberry.Std.AddOutputDriver(logberry.NewTextOutput(&TestingAdapter{t: t}, "testing"))
	logberry.Main = logberry.Std.Task("Test")

}
