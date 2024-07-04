package main

import (
	"strings"
	"testing"
)

func TestStore(t *testing.T) {
	t.Run("writeStream", func(t *testing.T) {
		t.Run("should return nil", func(t *testing.T) {
			s := NewStore(StoreOpts{})
			err := s.writeStream("key", strings.NewReader("data"))
			if err != nil {
				t.Errorf("expected nil, got %v", err)
			}
		})
	})
}
