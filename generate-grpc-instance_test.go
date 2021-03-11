package model_convert

import (
	"fmt"
	"testing"
)

func TestGenerateGRPCInstance(t *testing.T) {
	var src = `
service UserCoin {
    rpc GetUserInfo(GetUserInfoRequest) returns (UserInfo){}
}
`


	rs := GenerateGRPCInstance(src, GenerateGRPCInstanceArg{
		PbPackagePath: "shangraomajiang/control/user/userPb",
		PackageName:   "userControl",
	})

	fmt.Println(rs)
}

func TestGenerateGRPCInstanceV2(t *testing.T) {
	var si = ServiceItem{
		MethodName:   "GetUserChargeLog",
		RequestName:  "ProMiniGameGetUserChargeLogRequest",
		ResponseName: "ProMiniGameGetUserInfoResponse",
	}

	var si2 = ServiceItem{
		MethodName:   "GetUserInfo",
		RequestName:  "ProMiniGameGetUserInfoRequest",
		ResponseName: "ProMiniGameGetUserInfoResponse",
	}

	rpcServer := RpcServer{
		ServiceName:  "UserChargeLog",
		ServiceItems: []ServiceItem{si, si2},
	}
	rs := GenerateGRPCInstanceV2(
		rpcServer,
		"project/control/statistics/statisticsPb",
		"statisticsControl",
	)

	fmt.Println(rs)
}
