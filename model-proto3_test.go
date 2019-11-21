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
	type ChuangguanConfig struct {
		Id         int `gorm:"column:id;default:" json:"id" form:"id"`
		GameId     int `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		GameAreaId int `gorm:"column:game_area_id;default:" json:"game_area_id" form:"game_area_id"`

		StageNumber       int `gorm:"column:stage_number;default:" json:"stage_number" form:"stage_number"`
		InitScore         int `gorm:"column:init_score;default:" json:"init_score" form:"init_score"`
		StageScore        int `gorm:"column:stage_score;default:" json:"stage_score" form:"stage_score"`
		NeedGold          int `gorm:"column:need_gold;default:" json:"need_gold" form:"need_gold"`
		NeedWatchAdsTimes int `gorm:"column:need_watch_ads_times;default:" json:"need_watch_ads_times" form:"need_watch_ads_times"`
		NeedShareTimes    int `gorm:"column:need_share_times;default:" json:"need_share_times" form:"need_share_times"`
		EasterMaxTimes    int `gorm:"column:easter_max_times;default:" json:"easter_max_times" form:"easter_max_times"`
		DailyMaxPayChange int `gorm:"column:daily_max_pay_change;default:" json:"daily_max_pay_change" form:"daily_max_pay_change"`
	}
	ps, setM, setP := GoModelToProto2(ChuangguanConfig{}, map[string]string{
		"${pb_pkg_name}":    "gamePb",
		"${model_pkg_name}": "gameModel",
//		"${start_index}": "1",
	})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}
