package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func newStore() *Store {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}
func tearDown(t *testing.T, s *Store) {
	if err := s.clear(); err != nil {
		t.Error(err)
	}
}
func TestStoreDelete(t *testing.T) {
	s := newStore()
	tearDown(t, s)
}

func TestWrite(t *testing.T) {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "Iamhappy"
	data := bytes.NewReader([]byte("Hello There!"))
	if err := s.WriteStream(key, data); err != nil {
		t.Error(err)
	}
	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)

	fmt.Println(string(b))

	// s.delete(key)

}
