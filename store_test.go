package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpeg data"))
	s.writeStream("myspecialpicture", data)

	if err := s.writeStream("myspecialpicture", data); err != nil {
		t.Fatalf("failed to write stream: %v", err)
	}
}
