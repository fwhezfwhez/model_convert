package model_convert

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGoModelToProto3(t *testing.T) {
	type U struct {
		Username string
		Password string
		Age      int
		Id       int32
		Config   json.RawMessage
	}
	ps, setM, setP := GoModelToProto3(U{})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}
