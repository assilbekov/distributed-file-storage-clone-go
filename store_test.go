package main

import (
	"bytes"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

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

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "catsbestpictures"
	data := []byte("some jpeg data")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Fatalf("failed to write stream: %v", err)
	}

	if err := s.Delete(key); err != nil {
		t.Fatalf("failed to delete key: %v", err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "catsbestpictures"
	data := []byte("some jpeg data")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Fatalf("failed to write stream: %v", err)
	}

	if ok := s.Has(key); !ok {
		t.Fatalf("expected key %s to exist", key)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Fatalf("failed to read stream: %v", err)
	}

	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read all: %v", err)
	}

	if !bytes.Equal(b, data) {
		t.Fatalf("expected data %s, got %s", string(data), string(b))
	}

	s.Delete(key)
}
