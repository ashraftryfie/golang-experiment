package w3ctrace

import (
	"errors"
	"net/http"
)

var (
	ErrNotImplemented   = errors.New("exercise not implemented yet")
	ErrInvalidFormat    = errors.New("invalid traceparent header format")
	ErrInvalidTraceID   = errors.New("invalid trace ID: must be 32 non-zero hex chars")
	ErrInvalidParentID  = errors.New("invalid parent ID: must be 16 non-zero hex chars")
	ErrUnsupportedVer   = errors.New("unsupported traceparent version")
)

type TraceContext struct {
	Version  string
	TraceID  string
	ParentID string
	Flags    string
}

func ParseTraceparent(header string) (*TraceContext, error) {
	return nil, ErrNotImplemented
}

func (tc TraceContext) String() string {
	return ""
}

func GenerateNewTraceContext() TraceContext {
	return TraceContext{}
}

func InjectTraceparent(req *http.Request, tc TraceContext) {
}

func ExtractTraceparent(r *http.Request) (*TraceContext, bool) {
	return nil, false
}
