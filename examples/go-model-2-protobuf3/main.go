package main

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/model_convert"
)

func main() {
	type U struct {
		Username string
		Password string
		Age      int
		Id       int32
		Config   json.RawMessage
	}
	ps, setM, setP := model_convert.GoModelToProto3(U{})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}
