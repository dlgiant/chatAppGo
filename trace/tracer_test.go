package trace

import (
	"bytes"
	"io"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)

	if tracer == nil {
		t.Error("Tracer is nil")
	} else {
		tracer.Trace("Tracing is working.")
		if buf.String() != "Tracing is working.\n" {
			t.Errorf("Trace should not write '%s'.", buf.String())
		}
	}
}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
