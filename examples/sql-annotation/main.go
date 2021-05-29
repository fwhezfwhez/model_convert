package main

import (
	"fmt"
	"github.com/fwhezfwhez/model_convert"
)

func main() {
	sql := `
create table sql_generate_note(
   id serial primary key,                                  -- 自增id，主键
   updated_at timestamp with time zone default now(),      -- 更新于，
   created_at timestamp with time zone default now(),      -- 创建于
   
   -- 用户id
   user_id integer,
   
   -- 平台id
   -- 游戏id
   game_id integer,
   
   -- 包渠道1
   -- 包渠道2
   app_channel varchar -- 包渠道3
)
`

	rs := model_convert.GenerateNote(sql)

	fmt.Println(rs)
}