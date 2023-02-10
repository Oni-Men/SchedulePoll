package sliceutil

import (
	"testing"
)

func TestSlicePop(t *testing.T) {
	_, _, err := Pop([]string{})
	if err == nil {
		t.Fatalf("expect getting error. but error is nil")
	}

	pop, slice, _ := Pop([]string{
		"apple",
	})
	if *pop != "apple" {
		t.Fatalf("expected pop is apple. but we got %s", *pop)
	}

	if len(slice) != 0 {
		t.Fatalf("expect slice is empty. but acutual is %d", len(slice))
	}

}

func TestSlicePush(t *testing.T) {

}

func TestSliceUnpush(t *testing.T) {
	_, _, err := Unpush([]string{})
	if err == nil {
		t.Fatalf("expect getting error. but error is nil")
	}

	unpush, slice, _ := Unpush([]string{
		"apple",
		"banana",
		"cherry",
	})
	if *unpush != "cherry" {
		t.Fatalf("expected pop is cherry. but we got %s", *unpush)
	}

	if len(slice) != 2 {
		t.Fatalf("expect slice length is 2. but acutual is %d", len(slice))
	}
}
