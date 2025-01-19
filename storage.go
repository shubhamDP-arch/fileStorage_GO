package main

import (
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
		Original: hashStr,
	}

}
func (p PathKey) Filname() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Original)
}

type PathKey struct {
	PathName string
	Original string
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
func (s *Store) writeStream(key string, r io.Reader) error {
	pathkey := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathkey.PathName, os.ModePerm); err != nil {
		return err
	}
	pathAndFileName := pathkey.Filname()
	f, err := os.Create(pathAndFileName)
	if err != nil {
		return nil
	}
	n, err := io.Copy(f, r)

	if err != nil {
		return err
	}
	log.Printf("written %d bytes to disk %s", n, pathAndFileName)
	return nil
}
