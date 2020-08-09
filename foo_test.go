package main

import (
	"fmt"
	"testing"
)

func TestFoo(t *testing.T) {
	var foo Foo
	foo.A = 100
	foo.B = []string{"hahaha"}
	foo.C = map[int]string{
		1: "aaa",
		2: "bbb",
	}
	s := NewSerializer()
	err := s.SerializeFoo(&foo)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(s.Dump())
}
