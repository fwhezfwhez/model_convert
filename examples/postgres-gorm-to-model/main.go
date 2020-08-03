package main

import (
	"fmt"

	"github.com/fwhezfwhez/model_convert"
	// "model_convert"
)

/*
create table user_info(
	id serial primary key,
	user_id integer not null,
	open_id varchar,
	union_id varchar,
	game_id integer,
	user_name varchar not null default '',
	header_url varchar not null default '',
	sex int not null default 1, -- 1男，2女
	created_at timestamp without time zone not null default now(),
	last_login_time timestamp without time zone not null default now()
);
insert into user_info(user_id, open_id, game_id, user_name) values(10086, 'py', 78, 'Tommy');
insert into user_info(user_id, open_id, game_id, user_name) values(10087, 'jq', 78, 'Kitty');

*/
func main() {
	dataSouce := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "localhost", "5432", "postgres", "game", "disable", "123")
	tableName := "user_info"
	fmt.Println(model_convert.TableToStructWithTag(dataSouce, tableName))
}
