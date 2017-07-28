package trace

import (
	"io"
	"fmt"
)

type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

type nilTracer struct{}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func (t *nilTracer) Trace(a...interface{}) {}

func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

func Off() Tracer {
	return &nilTracer{}
}