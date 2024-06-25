package p2p

import "io"

type Decoder interface {
	Decode(io.Reader, any) error
}

type GOBDecoder struct{}

func (g *GOBDecoder) Decode(r io.Reader, v any) error {
	return nil
}
