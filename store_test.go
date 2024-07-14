package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestStoreDelete(t *testing.T) {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "MoiSpecial"
	data := []byte("i love her")

	if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if err := s.delete(key); err != nil {
		t.Error(err)
	}

}

func TestWrite(t *testing.T) {
	opts := StoreOps{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("Hello There!"))
	if err := s.WriteStream("Iamhappy", data); err != nil {
		t.Error(err)
	}
	r, err := s.Read("Iamhappy")
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)

	fmt.Println(string(b))

}
