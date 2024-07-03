package main

import "io"

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc
}

type Store struct {
	//
}

func (s *Store) writeStream(key string, r io.Reader) error {
	return nil
}
