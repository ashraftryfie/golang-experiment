package rot13

import (
	"io"
)

type Rot13Reader struct {
	r io.Reader
}

func NewRot13Reader(r io.Reader) *Rot13Reader {
	return &Rot13Reader{r: r}
}

func (rot *Rot13Reader) Read(p []byte) (n int, err error) {
	n, err = rot.r.Read(p)
	for i := 0; i < n; i++ {
		b := p[i]
		switch {
		case b >= 'a' && b <= 'z':
			p[i] = 'a' + (b-'a'+13)%26
		case b >= 'A' && b <= 'Z':
			p[i] = 'A' + (b-'A'+13)%26
		}
	}
	return n, err
}
