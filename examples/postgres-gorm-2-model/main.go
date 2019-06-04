package main

import (
	"fmt"
	"model_convert"
)

func main() {
	dataSouce := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "localhost", "5432", "postgres", "game", "disable", "123")
	tableName := "user_info"
	fmt.Println(model_convert.TableToStructWithTag(dataSouce, tableName))
}
