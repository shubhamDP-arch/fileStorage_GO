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
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformfunc,
	}
	s := NewStore(opts)
	key := "momsSpecial"
	data := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil{
		t.Error(err)
	}
	r, err := s.Read(key)
	if err != nil{
		t.Error(err)
	}
	b, _ := ioutil.ReadAll(r)
	fmt.Println(string(b))
	if string(b) != string(data) {
		t.Errorf("want %s hava %s",data, b)
	}
	
}
