package model_convert

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestGenerateMDAdd(t *testing.T) {
	type ResendRecord struct {
		Id         int             `gorm:"column:id;default:" json:"id" form:"id"`
		GameId     int             `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		GameAreaId int             `gorm:"column:game_area_id;default:" json:"game_area_id" form:"game_area_id"`
		UserId     int             `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		PropId     int             `gorm:"column:prop_id;default:" json:"prop_id" form:"prop_id"`
		PropNum    int             `gorm:"column:prop_num;default:" json:"prop_num" form:"prop_num"`
		ExpireIn   int             `gorm:"column:expire_in;default:" json:"expire_in" form:"expire_in"`
		Raw        json.RawMessage `gorm:"column:raw;default:" json:"raw" form:"raw"`
	}
	rs := generateMDAdd(ResendRecord{}, map[string]interface{}{
		"${model_chinese_name}": "补单记录",
		"${md_order}":           2,
	})
	fmt.Println(rs)
}

func TestGenerateMDUpdate(t *testing.T) {
	type JiugonggeSmallAwardPool struct {
		Id         int     `gorm:"column:id;default:" json:"id" form:"id"`
		GameId     int     `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		BlockId    int     `gorm:"column:block_id;default:" json:"block_id" form:"block_id"`
		PropId     int     `gorm:"column:prop_id;default:" json:"prop_id" form:"prop_id"`
		PropNum    int     `gorm:"column:prop_num;default:" json:"prop_num" form:"prop_num"`
		ExpireIn   int     `gorm:"column:expire_in;default:" json:"expire_in" form:"expire_in"`
		Percentage float64 `gorm:"column:percentage;default:" json:"percentage" form:"percentage"`
		TotalNum   int     `gorm:"column:total_num;default:" json:"total_num" form:"total_num"`
	}
	rs := generateMDUpdate(JiugonggeSmallAwardPool{}, map[string]interface{}{
		"${model_chinese_name}": "九宫格小奖池",
		"${md_order}":           1,
	})
	fmt.Println(rs)
}

func TestGenerateMDList(t *testing.T) {
	type ResendRecord struct {
		Id         int             `gorm:"column:id;default:" json:"id" form:"id"`
		GameId     int             `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		GameAreaId int             `gorm:"column:game_area_id;default:" json:"game_area_id" form:"game_area_id"`
		UserId     int             `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		PropId     int             `gorm:"column:prop_id;default:" json:"prop_id" form:"prop_id"`
		PropNum    int             `gorm:"column:prop_num;default:" json:"prop_num" form:"prop_num"`
		ExpireIn   int             `gorm:"column:expire_in;default:" json:"expire_in" form:"expire_in"`
		Raw        json.RawMessage `gorm:"column:raw;default:" json:"raw" form:"raw"`
	}
	rs := generateMDList(ResendRecord{}, map[string]interface{}{
		"${model_chinese_name}": "补发记录",
		"${md_order}":           2,
	})
	fmt.Println(rs)
}

func TestGenerateMDelete(t *testing.T) {
	type JiugonggeSmallAwardPool struct {
		Id         int     `gorm:"column:id;default:" json:"id" form:"id"`
		GameId     int     `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		BlockId    int     `gorm:"column:block_id;default:" json:"block_id" form:"block_id"`
		PropId     int     `gorm:"column:prop_id;default:" json:"prop_id" form:"prop_id"`
		PropNum    int     `gorm:"column:prop_num;default:" json:"prop_num" form:"prop_num"`
		ExpireIn   int     `gorm:"column:expire_in;default:" json:"expire_in" form:"expire_in"`
		Percentage float64 `gorm:"column:percentage;default:" json:"percentage" form:"percentage"`
		TotalNum   int     `gorm:"column:total_num;default:" json:"total_num" form:"total_num"`
	}
	rs := generateMDDelete(JiugonggeSmallAwardPool{}, map[string]interface{}{
		"${model_chinese_name}": "九宫格小奖池",
		"${md_order}":           1,
	})
	fmt.Println(rs)
}

func TestGenerateMGet(t *testing.T) {
	type JiugonggeSmallAwardPool struct {
		Id         int     `gorm:"column:id;default:" json:"id" form:"id"`
		GameId     int     `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		BlockId    int     `gorm:"column:block_id;default:" json:"block_id" form:"block_id"`
		PropId     int     `gorm:"column:prop_id;default:" json:"prop_id" form:"prop_id"`
		PropNum    int     `gorm:"column:prop_num;default:" json:"prop_num" form:"prop_num"`
		ExpireIn   int     `gorm:"column:expire_in;default:" json:"expire_in" form:"expire_in"`
		Percentage float64 `gorm:"column:percentage;default:" json:"percentage" form:"percentage"`
		TotalNum   int     `gorm:"column:total_num;default:" json:"total_num" form:"total_num"`
	}
	rs := generateMDGet(JiugonggeSmallAwardPool{}, map[string]interface{}{
		"${model_chinese_name}": "九宫格小奖池",
		"${md_order}":           1,
	})
	fmt.Println(rs)
}

func TestGenerateMD2(t *testing.T) {
	type FragmentGroup struct {
		Id        int       `gorm:"column:id;default:" json:"id" form:"id"`
		UpdatedAt time.Time `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`
		CreatedAt time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		GroupKey  string    `gorm:"column:group_key;default:" json:"group_key" form:"group_key"`
		GroupDesc string    `gorm:"column:group_desc;default:" json:"group_desc" form:"group_desc"`
	}
	rs := GenerateMarkDown(FragmentGroup{}, map[string]interface{}{
		"${model_chinese_name}": "碎片道具组",
		"${md_order}":           1,
	})
	fmt.Println(rs)
}
