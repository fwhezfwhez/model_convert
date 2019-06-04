package model_convert

import (
	"fmt"
	"testing"
)

func TestSplit2(t *testing.T) {
	tmp := ""
	rs := make([]string, 0, 20)
	Split2("   Name            string", " ", &tmp, &rs)
	fmt.Println(rs, tmp, len(rs))
}
func TestSplit(t *testing.T) {

	rs := Split(`
-- group:1
-- 1. 增加football_match表列 FFF
alter table football_match add column FFF varchar not null default '';

-- group:2
alter table football_match add column FFF2 varchar not null default '';
`, "\n")
	fmt.Println(rs, len(rs))
}

func TestAddJSONFormGormTag(t *testing.T) {
	rs := AddJSONFormGormTag(`
	type DbStruct struct{
		ColumnNumber int
		ColumnName string
		DataType string
	}
	`)
	fmt.Println(rs)
}

func TestHumpToUnderLine(t *testing.T) {
	fmt.Println(HumpToUnderLine("NameStructAge"))
}

func TestUnderLineToHump(t *testing.T) {
	fmt.Println(UnderLineToHump("NameStructAge"))
}

func TestFindUpperIndex(t *testing.T) {
	fmt.Println(FindUpperElement("NameStructAge"))
	var a = make(map[string]string, 0)
	a["1"] = a["1"] + "22"
	fmt.Println(a["1"])
}

//  generate model without model from database
func TestTableToStruct(t *testing.T) {
	dataSouce := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "localhost", "5432", "postgres", "test", "disable", "123")
	tableName := "football_match"
	fmt.Println(TableToStruct(dataSouce, tableName))
}


// generate model with json/form/gorm tag from database
func TestTableToStructWithTag(t *testing.T) {
	dataSouce := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "localhost", "5432", "postgres", "test", "disable", "123")
	tableName := "xyx_error"
	fmt.Println(TableToStructWithTag(dataSouce, tableName))
}


// add json and form for a go model
func TestAddJSONFormTag(t *testing.T) {
	fmt.Println(AddJSONFormTag(`
		type NotifyNormal struct {
	Url      string
	Data     json.RawMessage
	Interval string
}
	`))
}


func TestModelConvert_Generate(t *testing.T) {
	mc :=ModelConvert{}
	dataSouce := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "localhost", "5432", "postgres", "test", "disable", "123")
	mc.Init("postgres",dataSouce,false)
	mc.SetFlags(GORM | FORM| JSON)
	//s,e:=mc.Generate(nil)
	s,e :=mc.Generate([]string{"addresses"}...)
	if e!=nil{
		panic(e)
	}
	fmt.Println(s)
}

