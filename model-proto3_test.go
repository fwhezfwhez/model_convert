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
	type UserProp struct {
		Id                  int   `gorm:"column:id;default:" json:"id" form:"id"`
		GameId              int   `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		PlatformId          int   `gorm:"column:platform_id;default:" json:"platform_id" form:"platform_id"`
		UserId              int   `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		PropConfigId        int   `gorm:"column:prop_config_id;default:" json:"prop_config_id" form:"prop_config_id"`
		PropNum             int   `gorm:"column:prop_num;default:" json:"prop_num" form:"prop_num"`
		PropFrom            int   `gorm:"column:prop_from;default:" json:"prop_from" form:"prop_from"`
		PropSupportExchange int   `gorm:"column:prop_support_exchange;default:" json:"prop_support_exchange" form:"prop_support_exchange"`
		PropSupportGive     int   `gorm:"column:prop_support_give;default:" json:"prop_support_give" form:"prop_support_give"`
		PropValidityTime    int64 `gorm:"column:prop_validity_time;default:" json:"prop_validity_time" form:"prop_validity_time"`
		InUse               int   `gorm:"column:in_use;default:" json:"in_use" form:"in_use"`
	}
	ps, setM, setP := GoModelToProto2(UserProp{}, map[string]string{
		"${pb_pkg_name}":    "shopPb",
		"${model_pkg_name}": "shopModel",
//		"${start_index}": "1",
	})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}
