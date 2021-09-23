package model_convert

import (
	"fmt"
	"testing"
)

func TestGeneratePgNote(t *testing.T) {

	sql := `
create table union_config(
    id serial primary key,
    game_id integer not null default 0,
    union_id integer not null default 0,
    
    is_free integer NOT NULL DEFAULT 1, --竞技赛是否免费,1不免费、2免费
    free_desk_num integer not null default 1, -- 免费空桌子数
    empty_desk_num integer not null default 1, -- 客户端视角展示的最大空桌子数  
  
    created_at timestamp with time zone default now(),
    updated_at timestamp with time zone default now(),
  
    unique(game_id, union_id)
);
`

	rs := GenerateNote(sql)

	fmt.Println(rs)
}
