package trace_my

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {

	var buf bytes.Buffer
	tracer := New(&buf)

	if tracer == nil {
		t.Error("Return from New should not be nil")
	} else {
		tracer.Trace("Hello trace package")
		if buf.String() != "Hello trace package" {
			t.Errorf("Trace should not write '%s'.", buf.String())
		}

	}

	// t.Error("we haven't written out test yet")
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("something")
}
