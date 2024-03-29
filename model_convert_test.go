package model_convert

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
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

func TestFindMysqlClms(t *testing.T) {
	return
	r := findMysqlColumns("ft:123@/test?charset=utf8&parseTime=True&loc=Local&allowNativePasswords=true", "t_user")
	fmt.Println(r)
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
	dataSouce := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "localhost", "5432", "postgres", "game", "disable", "123")
	tableName := "user_migrate_process"
	fmt.Println(TableToStruct(dataSouce, tableName))
}

// generate model with json/form/gorm tag from database
func TestTableToStructWithTag(t *testing.T) {
	dataSouce := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "localhost", "5432", "postgres", "clubdb", "disable", "123456")
	tableName := "club_union"
	fmt.Println(TableToStructWithTag(dataSouce, tableName, map[string]interface{}{
		"dialect":            "postgres",
		"${db_instance}":     "db.QipaiDB",
		"${db_instance_pkg}": "configapisrv/dependency/db",
	}))
}

// add json and form for a go model
func TestAddJSONFormTag(t *testing.T) {
	fmt.Println(time.Now().Unix())
	fmt.Println(AddJSONFormTag(`
	`))
}

func TestModelConvert_Generate(t *testing.T) {
	mc := ModelConvert{}
	dataSouce := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "localhost", "5432", "postgres", "test", "disable", "123")
	mc.Init("postgres", dataSouce, false)
	mc.SetFlags(GORM | FORM | JSON)
	//s,e:=mc.Generate(nil)
	s, e := mc.Generate([]string{"addresses"}...)
	if e != nil {
		panic(e)
	}
	fmt.Println(s)
}

func TestTableToStructMysql(t *testing.T) {
	dataSouce := "ft:123@/test?charset=utf8&parseTime=True&loc=Local&allowNativePasswords=true"
	tableName := "t_user"
	fmt.Println(TableToStructWithTag(dataSouce, tableName, map[string]interface{}{
		"dialect":            "postgres",
		"${db_instance}":     "db.DB",
		"${db_instance_pkg}": "shangraomajiang/util/db",
	}))
}

func TestJSON(t *testing.T) {
	var b = `
{
    "buf": "eyJnYW1lX2lkIjo3OCwidXNlcl9pZCI6MzYwNjEwNjZ9"
}
`
	type B struct {
		Buf []byte `json:"buf"`
	}
	var bu B
	if e := json.Unmarshal([]byte(b), &bu); e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Println(string(bu.Buf))
}
