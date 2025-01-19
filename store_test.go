package main

import (
	"bytes"
	"testing"
)

func TestPathTrasformFunc(t *testing.T){
	
}
func TestStore(t *testing.T){
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)
	
	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("mySpecPicture", data); err!= nil{
		t.Error(err)
	}
}
