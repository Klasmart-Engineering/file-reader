package test

import (
	"file_reader/src"
	"testing"
)

func TestGetNumber(t *testing.T) {
	a := src.Foo()
	if a != 3 {
		t.Errorf("src.Foo() = %d; want 3", a)
	}
}
