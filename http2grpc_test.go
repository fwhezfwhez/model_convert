package model_convert

import (
	"fmt"
	"testing"
)

func TestServiceItem_String(t *testing.T) {
	var si = ServiceItem{
		MethodName:   "GetUserCoin",
		ResponseName: "MinigameGetUserCoinResponse",
		RequestName:  "MinigameGetUserCoinRequest",
	}

	fmt.Println(si.String())

}

func TestServiceItems_String(t *testing.T) {
	var si = ServiceItem{
		MethodName:   "GetUserCoin",
		ResponseName: "MinigameGetUserCoinResponse",
		RequestName:  "MinigameGetUserCoinRequest",
	}

	var si2 = ServiceItem{
		MethodName:   "GetUserCoin2",
		ResponseName: "MinigameGetUserCoinResponse2",
		RequestName:  "MinigameGetUserCoinRequest2",
	}

	var sis = sis([]ServiceItem{si, si2})
	fmt.Println(sis.String())
}

func TestRPCServer_String(t *testing.T) {
	var si = ServiceItem{
		MethodName:   "GetUserCoin",
		ResponseName: "MinigameGetUserCoinResponse",
		RequestName:  "MinigameGetUserCoinRequest",
	}

	var si2 = ServiceItem{
		MethodName:   "GetUserCoin2",
		ResponseName: "MinigameGetUserCoinResponse2",
		RequestName:  "MinigameGetUserCoinRequest2",
	}

	var rpcServer = RpcServer{
		ServiceName:  "coin",
		ServiceItems: []ServiceItem{si, si2},
	}

	fmt.Println(rpcServer.String())
}
