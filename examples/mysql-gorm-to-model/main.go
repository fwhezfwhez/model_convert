package main

import (
	"fmt"
)

func main() {
	dataSouce := "ft:123@/test?charset=utf8&parseTime=True&loc=Local&allowNativePasswords=true"
	tableName := "t_user"
	fmt.Println(TableToStructWithTag(dataSouce, tableName, "mysql"))
}
