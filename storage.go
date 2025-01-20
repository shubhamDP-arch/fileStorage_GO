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

func CASPathTransformfunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}

}
func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

type PathKey struct {
	PathName string
	FileName string
}
type PathTransformFunc func(string) PathKey
type Store struct {
	StoreOpts
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}
func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}
func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	f, err := os.Open(pathKey.FullPath())
	if err != nil {
		return nil, err
	}
	return f, err
}
func (s *Store) writeStream(key string, r io.Reader) error {
	pathkey := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathkey.PathName, os.ModePerm); err != nil {
		return err
	}
	FullPath := pathkey.FullPath()
	f, err := os.Create(FullPath)
	if err != nil {
		return nil
	}
	n, err := io.Copy(f, r)

	if err != nil {
		return err
	}
	log.Printf("written %d bytes to disk %s", n, FullPath)
	return nil
}
