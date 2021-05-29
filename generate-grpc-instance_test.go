package model_convert

import (
	"fmt"
	"testing"
)

func TestGenerateGRPCInstance(t *testing.T) {
	var src = `
	service ChangeMatchService{
		rpc GetUserChangeMatchStatus(ChaneMatchUserStatusReq) returns (ChaneMatchUserStatusResp){}
	}
`

	rs := GenerateGRPCInstance(src, GenerateGRPCInstanceArg{
		PbPackagePath: "path/to/challengePb",
		PackageName:   "packageaname",
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
