package userModel

import "github.com/fwhezfwhez/model_convert/examples/go-model-to-protobuf3/userProto"

func SetModelUserInfo(src userProto.UserInfo) UserInfo {
	var dest UserInfo
	dest.Arr = src.Arr
	dest.Config2 = src.Config2
	dest.Config = src.Config
	dest.Username = src.Username
	dest.Password = src.Password
	dest.Age = src.Age
	dest.Id = src.Id
	return dest
}


func SetProtoUserInfo(src UserInfo) userProto.UserInfo {
	var dest userProto.UserInfo
	dest.Arr = src.Arr
	dest.Config2 = src.Config2
	dest.Config = src.Config
	dest.Username = src.Username
	dest.Password = src.Password
	dest.Age = src.Age
	dest.Id = src.Id
	return dest
}
