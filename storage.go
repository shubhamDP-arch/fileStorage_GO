package main

import (
	"io"
	"log"
	"os"
)

type PathTransformFunc func(string) string
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
	pathName := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err
	}
	filename := "somefilename"
	pathAndFileName := pathName + "/" + filename
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
