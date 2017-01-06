package testing

import (
	"bytes"
	"testing"
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

func SetStdTesting(t *testing.T) {

	Std = NewRoot(24)
	Std.AddOutputDriver(NewTextOutput(&TestingAdapter{t: t}, "testing"))
	Main = &Task{
		uid:       newtaskuid(),
		component: "testing",
		activity:  "Component main",
		root:      Std,
	}

}
