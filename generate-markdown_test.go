package model_convert

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateMDAdd(t *testing.T) {
	type CashToHuafeiRecord struct {
		Id     int `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`

		Os         string `gorm:"column:os;default:" json:"os" form:"os"`
		AppChannel string `gorm:"column:app_channel;default:" json:"app_channel" form:"app_channel"`
		PlatformId int    `gorm:"column:platform_id;default:" json:"platform_id" form:"platform_id"`

		Phone    string `gorm:"column:phone;default:" json:"phone" form:"phone"`
		Realname string `gorm:"column:realname;default:" json:"realname" form:"realname"`

		OpType        string    `gorm:"column:op_type;default:" json:"op_type" form:"op_type"`
		Description   string    `gorm:"column:description;default:" json:"description" form:"description"`
		BalanceAmount int       `gorm:"column:balance_amount;default:" json:"balance_amount" form:"balance_amount"`
		HuafeiAmount  int       `gorm:"column:huafei_amount;default:" json:"huafei_amount" form:"huafei_amount"`
		HasHandled    int       `gorm:"column:has_handled;default:" json:"has_handled" form:"has_handled"`
		HandleAt      time.Time `gorm:"column:handle_at;default:" json:"handle_at" form:"handle_at"`
		CreatedAt     time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		UpdatedAt     time.Time `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`
	}
	rs := generateMDAdd(CashToHuafeiRecord{}, map[string]interface{}{
		"${model_chinese_name}": "余额提现话费记录",
		"${md_order}":           1,
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
	type DzpTaskUnion struct {
		Id            int    `gorm:"column:id;default:" json:"id" form:"id"`
		TaskUnionId   int    `gorm:"column:task_union_id;default:" json:"task_union_id" form:"task_union_id"`
		Description   string `gorm:"column:description;default:" json:"description" form:"description"`
		OsDescription string `gorm:"column:os_description;default:" json:"os_description" form:"os_description"`

		State     int    `gorm:"column:state;default:" json:"state" form:"state"`
		TaskTimes int    `gorm:"column:task_times;default:" json:"task_times" form:"task_times"`
		RefreshAt string `gorm:"column:refresh_at;default:" json:"refresh_at" form:"refresh_at"`
	}

	rs := generateMDList(DzpTaskUnion{}, map[string]interface{}{
		"${model_chinese_name}": "任务池",
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
	type ShopPropBackupRecord struct {
		Id          int       `gorm:"column:id;default:" json:"id" form:"id"`
		CreatedAt   time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		UpdatedAt   time.Time `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`
		ItemKey     string    `gorm:"column:item_key;default:" json:"item_key" form:"item_key"`
		Decription  string    `gorm:"column:decription;default:" json:"decription" form:"decription"`
		Title       string    `gorm:"column:title;default:" json:"title" form:"title"`
		DailyNum    int       `gorm:"column:daily_num;default:" json:"daily_num" form:"daily_num"`
		LifetimeNum int       `gorm:"column:lifetime_num;default:" json:"lifetime_num" form:"lifetime_num"`
	}
	rs := GenerateMarkDown(ShopPropBackupRecord{}, map[string]interface{}{
		"${model_chinese_name}": "库存信息",
		"${md_order}":           1,
	})
	fmt.Println(rs)
}
