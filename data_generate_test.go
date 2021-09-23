package model_convert

import (
	// "github.com/stretchr/testify/assert"
	"testing"
)

// 自己书写的前端接收请求结构体
type AppCoinTaxReq struct {
	StartTime    string `json:"start_date"`
	EndTime      string `json:"end_date"`
	GameId       int    `json:"game_id"`
	GameAreaId   int    `json:"game_area_id"`
	AppChannel   string `json:"app_channel"`
	LoginChannel string `json:"login_channel"`
	IsDownload   bool   `json:"is_download"`
}

// 自己书写的数据库结构体
type AppCoinTax struct {
	CalDate        string `json:"cal_date" gorm:"cal_date"`
	GameId         int    `json:"game_id" gorm:"game_id"`                   // 用户平台id
	GameAreaId     int    `json:"game_area_id" gorm:"game_area_id"`         // --子游戏id
	AppChannel     string `json:"app_channel" gorm:"app_channel"`           // --子渠道
	LoginChannel   string `json:"login_channel" gorm:"login_channel"`       // --登陆渠道
	GameRule       int    `json:"game_rule" gorm:"game_rule"`               // -- 场次级别
	AmountTax      int64  `json:"amount_tax" gorm:"amount_tax"`             // 平台金币税收
	PlatformGameId int    `json:"platform_game_id" gorm:"platform_game_id"` // 游戏平台id
}

func TestA(t *testing.T) {
	t.Run("测试", func(t *testing.T) {
		//确定传入参数
		needs := map[string]interface{}{
			//请求结构体
			"req_struct": AppCoinTaxReq{},
			//请求结构体注释(依次书写,如果不存在用"",而不是不填),如果不需要写注释可以省略 "notes_req_src"
			"notes_req_src": []string{"日期", "用户平台id", "--子游戏id", "--登陆渠道", "-- 场次级别"},
			//数据库结构体
			"struct": AppCoinTax{},
			//数据库结构体注释,同请求结构体注释
			"notes_src": []string{"日期", "用户平台id", "--子游戏id", "--登陆渠道", "-- 场次级别"},
			//绝对路径,会在该指定目录下创建文件夹和文件夹内的文件(以数据库结构体的名称命名文件夹)
			"path": "",
			//渠道汇总sql查询语句和其参数
			"sql_one": []string{"select * from table_name where id=? and name=?", "GameId", "EndTime"},
			//单渠道sql查询语句和其参数
			"sql_two": []string{"select * from table_name where id=? and name=?", "GameId", "EndTime"},
			//下载表格名
			"excel_table_name": "123",
		}
		e := GenerateData(needs)

		if e != nil {
			panic(e)
		}
		//该行可能会报错，可以添加 "github.com/stretchr/testify"
		//assert.Nil(t, e)
	})

}
