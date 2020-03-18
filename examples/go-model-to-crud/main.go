package main

import (
	"encoding/json"
	"fmt"

	mc "github.com/fwhezfwhez/model_convert"
)

func main() {
	type AUser struct {
		Id     int             `gorm:"column:id;default:" json:"id" form:"id"`
		Name   string          `gorm:"column:name;default:" json:"name" form:"name"`
		Age    int             `gorm:"column:age;default:" json:"age" form:"age"`
		Attach json.RawMessage `gorm:"column:attach;default:" json:"attach" form:"attach"`
	}

	rs := mc.GenerateCRUD(AUser{}, map[string]string{
		"${model}":        "userModel.AUser",
		"${handle_error}": `common.SaveError(errorx.Wrap(e))`,
		"${db_instance}":  "db.DB",
	})
	fmt.Println(rs)
}
