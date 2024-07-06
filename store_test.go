package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "catsbestpictures"
	pathname := CASPathTransformFunc(key)
	expectedPathname := "e24fc/4bc21/80e4d/f3696/836ab/8ccb8/ebe1b/7bf9b"
	if pathname != expectedPathname {
		t.Fatalf("expected pathname %s, got %s", expectedPathname, pathname)
	}
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
