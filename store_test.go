package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "catsbestpictures"
	pathname := CASPathTransformFunc(key)
	fmt.Printf("pathname: %s\n", pathname)
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
