package model_convert

import (
	"fmt"
	"strings"
)

// <service_item>
// Generate rpc item, like ` rpc GetUserCoin(MinigameGetUserCoinRequest) returns (MinigameGetUserCoinResponse) {}`
type ServiceItem struct {
	MethodName   string
	RequestName  string
	ResponseName string
}

func (si ServiceItem) String() string {
	rs := fmt.Sprintf(`
    rpc %s(%s) returns (%s) {}
`, si.MethodName, si.RequestName, si.ResponseName)

	return strings.Trim(rs, "\n")
}

// </service_item>

// <service_items>
// Generate items like
/*
    rpc GetUserCoin(MinigameGetUserCoinRequest) returns (MinigameGetUserCoinResponse) {}
    rpc GetUserCoin2(MinigameGetUserCoinRequest2) returns (MinigameGetUserCoinResponse2) {}
*/
type sis []ServiceItem

func (sis sis) String() string {
	var rs = make([]string, 0, 10)
	for _, v := range sis {
		rs = append(rs, v.String())
	}

	return strings.Join(rs, "\n")
}

// </service_items>

//<rpc_server>
// Generate rpc server node like:
/*
service coin {
    rpc GetUserCoin(MinigameGetUserCoinRequest) returns (MinigameGetUserCoinResponse) {}
    rpc GetUserCoin2(MinigameGetUserCoinRequest2) returns (MinigameGetUserCoinResponse2) {}
}
*/
type RpcServer struct {
	ServiceName  string        // rpc service名
	ServiceItems []ServiceItem // rpc 方法行
}

func (rs RpcServer) String() string {
	return fmt.Sprintf(`
service %s{
%s
}
`, rs.ServiceName, sis(rs.ServiceItems).String())
}

//</rpc_server>

// 只支持POST http请求，转换为grpc, request只支持非内嵌结构
// ${methodName}
// ${gin_router}
// ${message_name}
func ServcieToGRPCService(request interface{}, replacement map[string]string) string {
	var rs string

	type CountResult struct {
		MessageName string // message名
	}

	//var protoRPC string
	//var goGRPC string
	//
	protoRequest, _, _ := GoModelToProto2(request, replacement)

	//GoModelToProto2()
	_ = protoRequest

	return rs
}

func handleHTTP2GRPCReplacement(replacement map[string]interface{}) {

}
