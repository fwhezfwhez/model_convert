package model_convert

import (
	"fmt"
	"testing"
)

func TestLowerFistLetter(t *testing.T) {
	fmt.Println(LowerFistLetter("Hello"))
	fmt.Println(LowerFistLetter("hello"))
	fmt.Println(LowerFistLetter(""))
	fmt.Println(LowerFistLetter("H"))
	fmt.Println(LowerFistLetter("h"))
}
