package model_convert

import (
	"fmt"
	"strings"
	"testing"
)

func TestGeneratePgNote(t *testing.T) {
	var line = "//abcdefghijk"
	fmt.Println(line[strings.Index(line, "//")+len("//"):])

	sql := `
create table app_share_user_process(
   id serial primary key,
   created_at timestamp with time zone default now(),
   updated_at timestamp with time zone default now(),
  
   app_channel varchar,      -- 包渠道
   game_id integer,          -- 平台id
   user_id integer,          -- 用户id
   share_key varchar,        -- 分享活动key
   
   qrcode_url varchar -- 二维码链接
)
`

	rs := GenerateNote(sql)

	fmt.Println(rs)
}
