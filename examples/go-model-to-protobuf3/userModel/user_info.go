package userModel

import (
	"encoding/json"
)

type ArrayElement struct {
}
type UserInfo struct {
	Arr     []ArrayElement
	Config2 []byte
	Config  json.RawMessage

	Username string
	Password string
	Age      int
	Id       int32
}
