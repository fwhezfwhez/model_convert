package main

import (
	"fmt"
	"model_convert"
	"model_convert/examples/go-model-to-protobuf3/userModel"
)

func main() {
	ps, setM, setP := model_convert.GoModelToProto3(userModel.UserInfo{}, map[string]string{
		"${pb_pkg_name}":    "userProto",
		"${model_pkg_name}": "empty",
	})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}
