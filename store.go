package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultPath = "Data"

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Filename string
}
type StoreOps struct {
	root              string
	PathTransformFunc PathTransformFunc
}
type Store struct {
	StoreOps
}

func (p PathKey) FullName() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashstr := hex.EncodeToString(hash[:])
	block := 5
	slice := len(hashstr) / block
	paths := make([]string, slice)
	for i := 0; i < block; i++ {
		from, to := i*block, (i*block)+block
		paths[i] = hashstr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashstr,
	}

}
func DefaultPathTransformFunc(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

func NewStore(opts StoreOps) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.root) == 0 {
		opts.root = defaultPath
	}
	return &Store{
		StoreOps: opts,
	}
}

func (s *Store) delete(key string) error {
	pathkey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("Deleted [%s] from disk", pathkey.Filename)
	}()
	return os.RemoveAll(strings.Split(pathkey.Pathname, "/")[0])
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.ReadStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) ReadStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)

	return os.Open(pathKey.FullName())
}
func (s *Store) WriteStream(key string, r io.Reader) error {
	PathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(PathKey.Pathname, os.ModePerm); err != nil {
		return err
	}
	fullPath := PathKey.FullName()
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer func() {
		f.Close()
	}()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("Succesfully Writen [%d] Bytes to disk at %s", n, fullPath)

	return nil
}
