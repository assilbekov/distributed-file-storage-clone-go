package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "catsbestpictures"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalKey := "e24fc4bc2180e4df3696836ab8ccb8ebe1b7bf9b"
	expectedPathname := "e24fc/4bc21/80e4d/f3696/836ab/8ccb8/ebe1b/7bf9b"
	if pathKey.PathName != expectedPathname {
		t.Fatalf("expected pathname %s, got %s", expectedPathname, pathKey.PathName)
	}

	if pathKey.Filename != expectedOriginalKey {
		t.Fatalf("expected original key %s, got %s", expectedOriginalKey, pathKey.Filename)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpeg data"))
	s.writeStream("myspecialpicture", data)

	if err := s.writeStream("myspecialpicture", data); err != nil {
		t.Fatalf("failed to write stream: %v", err)
	}
}
