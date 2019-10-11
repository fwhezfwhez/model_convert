package model_convert

import (
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
	fmt.Println(GenerateWhere(VxTemplateUser{}, false, map[string]string{
		"${jump_fields}": "password,game_id",
	}))
}
