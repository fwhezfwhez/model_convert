package model_convert

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestGenerateWhere(t *testing.T) {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string    `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int       `gorm:"column:state;default:" json:"state" form:"state"`
		CreatedAt  time.Time `gorm:"column:created_at;default:" json:"created_at"`
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
		"${model}":        "payModel.LingqianOrder",
		"${handler_name}": "HTTPListLingqianOrder",
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

	rs := GenerateGetOneAPI(VxTemplateUser{}, map[string]string{
		"${model}":        "payModel.LingqianOrder",
		"${handler_name}": "HTTPGetOneLingqianOrder",
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

	rs := GenerateAddOneAPI(VxTemplateUser{}, map[string]string{
		"${model}":        "payModel.LingqianOrder",
		"${handler_name}": "HTTPAddLingqianOrder",
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
		"${model}":        "payModel.LingqianOrder",
		"${handler_name}": "HTTPDeleteLingqianOrder",
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
		"${model}":              "payModel.LingqianOrder",
		"${handler_name}":       "HTTPDeleteLingqianOrder",
		"${handle_error}":       `common.SaveError(e)`,
		"${args_forbid_update}": "UserId, game_id",
	})
	fmt.Println(rs)
}

// TODO: 增加缓存CRUD支持
func TestGenerateCRUD(t *testing.T) {
	type CardLibrary struct {
		Id         int             `gorm:"column:id;default:" json:"id" form:"id"`
		Type       int             `gorm:"column:type;default:" json:"type" form:"type"`
		Infomation string          `gorm:"column:infomation;default:" json:"infomation" form:"infomation"`
		Cards      json.RawMessage `gorm:"column:cards;default:" json:"cards" form:"cards"`
		Stat       bool            `gorm:"column:stat;default:" json:"stat" form:"stat"`
		CreatedAt  time.Time       `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		ModifyAt   time.Time       `gorm:"column:modify_at;default:" json:"modify_at" form:"modify_at"`
		IsDelete   bool            `gorm:"column:is_delete;default:" json:"is_delete" form:"is_delete"`
	}
	rs := GenerateCRUD(CardLibrary{}, map[string]string{
		"${generate_to_pkg}": "model",
		"${model}":           "model.CardLibrary",
		"${handle_error}":    "fmt.Println(errorx.Wrap(e))",
		"${db_instance}":     "db.DB",
	})
	_ = rs
	fmt.Println(rs)
}
