package rot13

import (
	"errors"
	"io"
)

var ErrNotImplemented = errors.New("TODO: implement")

type Rot13Reader struct {
	r io.Reader
}

func NewRot13Reader(r io.Reader) *Rot13Reader {
	return &Rot13Reader{r: r}
}

func (rot *Rot13Reader) Read(p []byte) (n int, err error) {
	// TODO: Call rot.r.Read(p), transform read bytes in-place, return n, err
	return 0, ErrNotImplemented
}
