package model_convert

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenerateWhere(t *testing.T) {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int    `gorm:"column:state;default:" json:"state" form:"state"`
	}
	fmt.Println(GenerateListWhere(VxTemplateUser{}, false, map[string]string{
		"${jump_fields}": "password,game_id",
	}))
}

func TestGenerateList(t *testing.T) {
	type LingqianOrder struct {
		Id       int             `gorm:"column:id;default:" json:"id" form:"id"`
		OrderId  string          `gorm:"column:order_id;default:" json:"order_id" form:"order_id"`
		GameId   int             `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId   int             `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId   string          `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`
		State    int             `gorm:"column:state;default:" json:"state" form:"state"`
		Request  json.RawMessage `gorm:"column:request;default:" json:"request" form:"request"`
		Response string          `gorm:"column:response;default:" json:"response" form:"response"`
	}

	fmt.Println(GenerateListAPI(LingqianOrder{}, false, map[string]string{
		"${model}": "payModel.LingqianOrder",
		"${handler_name}" : "HTTPListLingqianOrder",
		"${handle_error}": `common.SaveError(e)`,
	}))
}

func TestGenerateGetOneAPI(t *testing.T) {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int    `gorm:"column:state;default:" json:"state" form:"state"`
	}

	rs :=GenerateGetOneAPI(VxTemplateUser{}, map[string]string{
		"${model}": "payModel.LingqianOrder",
		"${handler_name}" : "HTTPGetOneLingqianOrder",
		"${handle_error}": `common.SaveError(e)`,
	})
	fmt.Println(rs)
}


func TestGenerateAddOneAPI(t *testing.T) {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int    `gorm:"column:state;default:" json:"state" form:"state"`
	}

	rs :=GenerateAddOneAPI(VxTemplateUser{}, map[string]string{
		"${model}": "payModel.LingqianOrder",
		"${handler_name}" : "HTTPAddLingqianOrder",
		"${handle_error}": `common.SaveError(e)`,
	})
	fmt.Println(rs)
}

func TestGenerateDeleteOneAPI(t *testing.T) {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int    `gorm:"column:state;default:" json:"state" form:"state"`
	}

	rs := GenerateDeleteOneAPI(VxTemplateUser{}, map[string]string{
		"${model}": "payModel.LingqianOrder",
		"${handler_name}" : "HTTPDeleteLingqianOrder",
		"${handle_error}": `common.SaveError(e)`,
	})
	fmt.Println(rs)
}

func TestGenerateUpdateOneAPI(t *testing.T) {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int    `gorm:"column:state;default:" json:"state" form:"state"`
	}

	rs := GenerateUpdateOneAPI(VxTemplateUser{}, map[string]string{
		"${model}": "payModel.LingqianOrder",
		"${handler_name}" : "HTTPDeleteLingqianOrder",
		"${handle_error}": `common.SaveError(e)`,
		"${args_forbid_update}": "UserId, game_id",
	})
	fmt.Println(rs)
}

// TODO: 增加缓存CRUD支持
func TestGenerateCRUD(t *testing.T) {
	type LingqianBalance struct {
		GameId        int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		GameAreaId    int    `gorm:"column:game_area_id;default:" json:"game_area_id" form:"game_area_id"`
		UserId        int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId        string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`
		Balance       int    `gorm:"column:balance;default:" json:"balance" form:"balance"`
		HistoryCashed int    `gorm:"column:history_cashed;default:" json:"history_cashed" form:"history_cashed"`
		Id            int    `gorm:"column:id;default:" json:"id" form:"id"`
		Origin        int    `gorm:"column:origin;default:" json:"origin" form:"origin"`
	}

	rs := GenerateCRUD(LingqianBalance{}, map[string]string{
		"${model}": "payModel.LingqianBalance",
		"${handle_error}": `common.SaveError(errorx.Wrap(e))`,
		"${db_instance}": "db.DB",
	})
	fmt.Println(rs)
}
