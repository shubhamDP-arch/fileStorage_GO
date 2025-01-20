package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)
const defaultRootFolderName  = "goodgame"
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
func (p PathKey)FirstPathName()string  {
	paths:= strings.Split(p.PathName, "/")[0]
	if len(paths) == 0{
		fmt.Errorf("%s the file path is too short", paths)
		return ""
	}
	return paths
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

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

type StoreOpts struct {
	Root 	string
	PathTransformFunc PathTransformFunc
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
		
	}
	if len(opts.Root) == 0{
		opts.Root = defaultRootFolderName
	}
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
func (s *Store)Has(key string)bool  {
	pathKey := s.PathTransformFunc(key)
	_, err:= os.Stat(pathKey.FullPath())
	return err == fs.ErrNotExist
}
func (s *Store)Delete(key string)error  {
	pathKey := s.PathTransformFunc(key)
	defer func ()  {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()
	// err := os.RemoveAll(pathKey.FullPath())
	// if err != nil {
	// 	return err
	// }
	// return nil
	return os.RemoveAll(pathKey.FirstPathName())
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
	if err := os.MkdirAll(s.Root+"/"+pathkey.PathName, os.ModePerm); err != nil {
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
