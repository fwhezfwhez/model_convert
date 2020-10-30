package model_convert

import (
	"fmt"
	"testing"
)

func TestGenerateGRPCInstance(t *testing.T) {
	var src = `
service UserChargeLog{
    rpc GetUserChargeLog(ProMiniGameGetUserChargeLogRequest) returns (ProMiniGameGetUserChargeLogResponse) {}
    rpc GetUserInfo(ProMiniGameGetUserInfoRequest) returns (ProMiniGameGetUserInfoResponse) {}
    rpc GetUserCoin(ProMiniGameGetUserCoinRequest) returns (ProMiniGameGetUserCoinResponse) {}
}
`
	rs := GenerateGRPCInstance(src, GenerateGRPCInstanceArg{
		PbPackagePath: "project/control/statistics/statisticsPb",
		PackageName:   "statisticsControl",
	})

	fmt.Println(rs)
}
