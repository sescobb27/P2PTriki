package triki

import (
	"testing"
)

func assertNoError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestMakeBoard(t *testing.T) {
	NewBoard()
}
