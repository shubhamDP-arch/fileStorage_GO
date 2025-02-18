package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPathTrasformFunc(t *testing.T) {

}
func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)
	key := "momsspecial"
	data := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if ok := s.Has(key); !ok {
		t.Errorf("expected to have key %s", key)
	}
	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := ioutil.ReadAll(r)
	fmt.Println(string(b))
	if string(b) != string(data) {
		t.Errorf("want %s hava %s", data, b)
	}

}

func TestStoreDeleteKey(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformfunc,
	}
	s := NewStore(opts)
	key := "foobar"
	data := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func newStore()*Store  {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformfunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store){
	if err := s.Clear(); err != nil{
		t.Error(err)
	}
}