package model_convert

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGoModelToProto3(t *testing.T) {
	type ArrayElement struct {
	}
	type U struct {
		Arr     []ArrayElement
		Config2 []byte
		Config  json.RawMessage

		Username string
		Password string
		Age      int
		Id       int32
	}
	ps, setM, setP := GoModelToProto3(U{}, map[string]string{
		"${pb_pkg_name}":    "userProto",
		"${model_pkg_name}": "userModel",
	})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}

func TestGoModelToProto2(t *testing.T) {
	type MailSendTaskNewReq struct {
		GameId         int     `json:"game_id" binding:"required"`
		MailBizType    int     `json:"mail_biz_type" binding:"required"`
		Title          string  `json:"title" binding:"required"`
		Content        string  `json:"content" binding:"required"`
		ReceiverType   int     `json:"receiver_type"  binding:"required"`
		ReceiverIds    []int64 `json:"receiver_ids" binding:"required"`
		SenderUserId   int     `json:"sender_user_id" binding:"required"`
		AttachmentFlag int     `json:"attachment_flag" binding:"required"`
		//AttachmentJson []Attachment `json:"attachment_json" binding:"required"`
		NotCreateTask bool `json:"not_create_task"` // 不创建任务，否代表默认创建任务，是代表不创建任务
	}
	ps, setM, setP := GoModelToProto2(MailSendTaskNewReq{}, map[string]string{
		"${pb_pkg_name}":    "userPb",
		"${model_pkg_name}": "userModel",
		//		"${start_index}": "1",
	})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}
