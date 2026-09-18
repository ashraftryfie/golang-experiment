package w3ctrace

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrInvalidFormat   = errors.New("invalid traceparent header format")
	ErrInvalidTraceID  = errors.New("invalid trace ID: must be 32 non-zero hex chars")
	ErrInvalidParentID = errors.New("invalid parent ID: must be 16 non-zero hex chars")
	ErrUnsupportedVer  = errors.New("unsupported traceparent version")
)

type TraceContext struct {
	Version  string
	TraceID  string
	ParentID string
	Flags    string
}

func isHex(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

func ParseTraceparent(header string) (*TraceContext, error) {
	parts := strings.Split(header, "-")
	if len(parts) != 4 {
		return nil, ErrInvalidFormat
	}

	version := parts[0]
	traceID := parts[1]
	parentID := parts[2]
	flags := parts[3]

	if version != "00" {
		return nil, ErrUnsupportedVer
	}

	if len(traceID) != 32 || !isHex(traceID) || traceID == strings.Repeat("0", 32) {
		return nil, ErrInvalidTraceID
	}

	if len(parentID) != 16 || !isHex(parentID) || parentID == strings.Repeat("0", 16) {
		return nil, ErrInvalidParentID
	}

	if len(flags) != 2 || !isHex(flags) {
		return nil, ErrInvalidFormat
	}

	return &TraceContext{
		Version:  version,
		TraceID:  traceID,
		ParentID: parentID,
		Flags:    flags,
	}, nil
}

func (tc TraceContext) String() string {
	return fmt.Sprintf("%s-%s-%s-%s", tc.Version, tc.TraceID, tc.ParentID, tc.Flags)
}

func GenerateNewTraceContext() TraceContext {
	traceBytes := make([]byte, 16)
	_, _ = rand.Read(traceBytes)

	parentBytes := make([]byte, 8)
	_, _ = rand.Read(parentBytes)

	return TraceContext{
		Version:  "00",
		TraceID:  hex.EncodeToString(traceBytes),
		ParentID: hex.EncodeToString(parentBytes),
		Flags:    "01",
	}
}

func InjectTraceparent(req *http.Request, tc TraceContext) {
	req.Header.Set("traceparent", tc.String())
}

func ExtractTraceparent(r *http.Request) (*TraceContext, bool) {
	raw := r.Header.Get("traceparent")
	if raw == "" {
		return nil, false
	}
	tc, err := ParseTraceparent(raw)
	if err != nil {
		return nil, false
	}
	return tc, true
}
