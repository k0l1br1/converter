package bins

import (
	"reflect"
	"testing"
)

const testBinsFile = "/tmp/test-bins.bin"

var bins = [...]Bin{
	{1, 1, 0.1},
	{0, 2, 0.2},
	{1, 3, 0.8},
}

func TestBinsNew(t *testing.T) {
	if err := InitNew(testBinsFile); err != nil {
		t.Fatal("init new bins file", err)
	}

	bs1 := bins[:2]
	if err := Write(testBinsFile, bs1); err != nil {
		t.Fatal("write to bins file", err)
	}

	bs2, err := ReadAll(testBinsFile)
	if err != nil {
		t.Fatal("read from bins file", err)
	}

	if !reflect.DeepEqual(bs1, bs2) {
		t.Fatalf("invalid read result: want %v, got %v", bs1, bs2)
	}
}

func TestBinsAppend(t *testing.T) {
    if err := Write(testBinsFile, bins[2:]); err != nil {
		t.Fatal("append to bins file", err)
	}

	bs2, err := ReadAll(testBinsFile)
	if err != nil {
		t.Fatal("read from bins file", err)
	}

    bs1 := bins[:]
	if !reflect.DeepEqual(bs1, bs2) {
		t.Fatalf("invalid read appended result: want %v, got %v", bs1, bs2)
	}
}
