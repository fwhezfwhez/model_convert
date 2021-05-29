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
	type PropTagConfig struct {
		Id        int       `gorm:"column:id;default:" json:"id" form:"id"`
		UpdatedAt time.Time `gorm:"column:updated_at;default:" json:"updated_at" form:"updated_at"`
		CreatedAt time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		PropTag   string    `gorm:"column:prop_tag;default:" json:"prop_tag" form:"prop_tag"`
		GameId    int       `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		Title     string    `gorm:"column:title;default:" json:"title" form:"title"`
	}

	rs := GenerateCRUD(PropTagConfig{}, map[string]string{
		"${generate_to_pkg}": "shopControl",
		"${model}":           "shopModel.PropTagConfig",
		"${handle_error}":    "shopUtil.SaveError(errorx.Wrap(e))",
		"${db_instance}":     "(shopModel.PropTagConfig{}).DB()",
	})
	_ = rs
	fmt.Println(rs)
}
