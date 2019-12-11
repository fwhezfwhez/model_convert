package main

import (
	"fmt"
	mc "github.com/fwhezfwhez/model_convert"
)

func main() {
	dataSouce := "ft:123@/test?charset=utf8&parseTime=True&loc=Local&allowNativePasswords=true"
	tableName := "t_user"
	fmt.Println(mc.TableToStructWithTag(dataSouce, tableName, "mysql"))
}
