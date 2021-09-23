package model_convert

import (
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"os"
	"path"
	"reflect"
	"strings"
)

type ReplaceType struct {
	//请求结构体
	ReqSrc interface{}
	//数据库结构体
	Src interface{}

	//路径
	Path string
	//文件夹名称
	NameDir string
	//调用方法数据表名单
	ExcelTableName string

	// 打印方法
	LogPrint string

	SrcInfo
	ReqSrcInfo
	SqlInfo
}

// SrcInfo 数据库结构体信息
type SrcInfo struct {
	//数据库结构体名称
	NameSrc string
	//数据库结构体注释
	NotesSrc []string
}

// ReqSrcInfo 请求结构体
type ReqSrcInfo struct {
	//请求结构体名称
	NameReqSrc string
	//请求结构体注释
	NotesReqSrc []string
}

type SqlInfo struct {
	//自定义表名
	NameTable string
	//第一个sql语句和其所需要的参数
	SqlOne []string
	//第二个sql语句和其所需要的参数
	SqlTwo []string
}

func New(needs map[string]interface{}) (*ReplaceType, error) {
	//确定needs完整
	checkParam(needs)

	//获取src中的名称
	//数据库结构体名称
	nameSrc := reflect.TypeOf(needs["struct"]).Name()
	//请求结构体名称
	nameReqSrc := reflect.TypeOf(needs["req_struct"]).Name()
	//格式化文件夹名称
	nameDir := strings.ToLower(string(nameSrc[0])) + nameSrc[1:]
	//自定义表名
	nameTable := addAndTransferToLower(nameDir)

	return &ReplaceType{
		LogPrint: needs["log_print"].(string),
		ReqSrc:   needs["req_struct"],
		Src:      needs["struct"],

		Path: needs["path"].(string),
		//文件夹名称
		NameDir: nameDir,

		ExcelTableName: needs["excel_table_name"].(string),

		SrcInfo: SrcInfo{
			NameSrc:  nameSrc,
			NotesSrc: needs["notes_src"].([]string),
		},
		ReqSrcInfo: ReqSrcInfo{
			NameReqSrc:  nameReqSrc,
			NotesReqSrc: needs["notes_req_src"].([]string),
		},
		SqlInfo: SqlInfo{
			NameTable: nameTable,
			SqlOne:    needs["sql_one"].([]string),
			SqlTwo:    needs["sql_two"].([]string),
		},
	}, nil
}

// checkParam 判断数据是否存在
func checkParam(needs map[string]interface{}) {
	if _, ok := needs["path"]; !ok {
		needs["path"], _ = os.Getwd()
	}
	if _, ok := needs["sql_one"]; !ok {
		needs["sql_one"] = []string{""}
	}
	if _, ok := needs["sql_two"]; !ok {
		needs["sql_two"] = []string{""}
	}
	if _, ok := needs["excel_table_name"]; !ok {
		needs["excel_table_name"] = ""
	}
	if _, ok := needs["notes_req_src"]; !ok {
		needs["notes_req_src"] = []string{""}
	}
	if _, ok := needs["notes_src"]; !ok {
		needs["notes_src"] = []string{""}
	}
	type s struct{}
	if _, ok := needs["req_struct"]; !ok {
		needs["req_struct"] = s{}
	}
	if _, ok := needs["struct"]; !ok {
		needs["struct"] = s{}
	}

	if _, ok := needs["log_print"]; !ok {
		needs["log_print"] = "fmt.Println"
	}
}

//  isParam 判断reqSrc中Param是否存在
func isParam(reqSrc interface{}, paramName string) (result bool) {
	typeOf := reflect.TypeOf(reqSrc)
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Name == paramName {
			result = true
			break
		}
	}
	return result
}

// GenerateData 自动创建大数据类型 (path：绝对参数,reqSrc:请求结构体，src:数据库结构体)
// 1.创建文件夹,文件
// 2.打开文件写入内容
// 3.格式化文件
func GenerateData(needs map[string]interface{}) error {

	replaceType, e := New(needs)
	if e != nil {
		return errorx.Wrap(e)
	}

	if e = replaceType.create(); e != nil {
		return errorx.Wrap(e)
	}

	if e = replaceType.write(); e != nil {
		return errorx.Wrap(e)
	}

	return nil
}

// 通过结构体解析
func (r ReplaceType) GenerateData() error {
	if e := r.create(); e != nil {
		return e
	}

	if e := r.write(); e != nil {
		return e
	}

	return nil
}

// create
// 创建文件夹，创建文件
func (r ReplaceType) create() error {
	//对名称进行格式化
	pathDir := path.Join(r.Path, r.NameDir)

	//创建文件夹
	if e := os.MkdirAll(pathDir, os.ModePerm); e != nil {
		return errorx.Wrap(e)
	}

	//创建文件
	creFile := []string{"router.go", "service.go", "model.go"}
	for _, v := range creFile {
		_, e := os.Create(path.Join(pathDir, v))
		if e != nil {
			return errorx.Wrap(e)
		}
	}
	return nil
}

// write 写入内容
//
func (r ReplaceType) write() error {
	if e := r.writeModel(); e != nil {
		return errorx.Wrap(e)
	}
	if e := r.writeService(); e != nil {
		return errorx.Wrap(e)
	}
	if e := r.writeRouter(); e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

// writeModel 写入model
// ${name_dir} 包名
// ${module_req} 根据请求结构体自动生成
// ${struct} 结构体
// ${struct_name} 结构体名称
// ${module_req_name} module结构体名称
// ${check_login_channel} 是否存在login_channel
func (r ReplaceType) writeModel() error {

	var front = `package ${name_dir}

import "path/to/db"

//根据请求结构体生成
${module_req}

// 数据库表结构体
${struct}


// 开发人员自己写sql
func Get${struct_name}QueryList(model ${module_req_name}) ([]${struct_name}, int, error) {
	a := make([]${struct_name}, 0)
	var count int

	${check_login_channel}
}
`

	front = strings.ReplaceAll(front, "${name_dir}", r.NameDir)
	front = strings.ReplaceAll(front, "${struct}", formatPrint(r.Src, r.NotesSrc))
	front = strings.ReplaceAll(front, "${struct_name}", r.NameSrc)
	front = strings.ReplaceAll(front, "${req_struct_name}", r.NameReqSrc)

	tmp := formatPrint(r.ReqSrc, r.NotesReqSrc)
	tmp = strings.ReplaceAll(tmp, r.NameReqSrc, r.NameSrc+"ModuleReq")
	front = strings.ReplaceAll(front, "${module_req}", tmp)
	front = strings.ReplaceAll(front, "${module_req_name}", r.NameSrc+"ModuleReq")
	front = strings.ReplaceAll(front, "${check_login_channel}", r.checkLoginChannel(isParam(r.ReqSrc, "LoginChannel")))

	if e := openAndWriteStringFile(path.Join(r.Path, r.NameDir, "model.go"), front); e != nil {
		return e
	}
	return nil
}

// checkLoginChannel 根据是否存在LoginChannel生成不同的结果
// ${sql_one} sqlOne需要参数
// ${sql_two} sqlTwo需要参数
func (r ReplaceType) checkLoginChannel(haslc bool) string {
	if haslc {
		var front = `
	if model.LoginChannel == "all" {
		${sql_one}
	} else {
		${sql_two}
	}
`
		front = strings.ReplaceAll(front, "${sql_one}", addContent(r.SqlInfo.SqlOne[0], r.SqlInfo.NameTable, r.SqlInfo.SqlOne[1:], r.ReqSrc, true))
		front = strings.ReplaceAll(front, "${sql_two}", addContent(r.SqlInfo.SqlTwo[0], r.SqlInfo.NameTable, r.SqlInfo.SqlTwo[1:], r.ReqSrc, false))
		return front
	} else {
		return addContent(r.SqlInfo.SqlOne[0], r.SqlInfo.NameTable, r.SqlInfo.SqlOne[1:], r.ReqSrc, false)
	}
}

// addContent 添加内部内容
// ${sql} 语句
// ${sql_params} 参数
// ${table_name} 表名
// ${middle} 中间数据
func addContent(sql, tableName string, sqlParams []string, src interface{}, hasGroup bool) string {
	var front = `
	sql := "${sql}"
	if err := db.DataDB.Table("${table_name}").Raw(sql,${sql_params}).Find(&a).Error; err != nil {
		return a, count, err
	}

	if err := db.DataDB.Table("${table_name}").
		Where("cal_date between ? and ?", model.StartTime, model.EndTime).
${middle}
		return a, count, err
	}
	return a, count, nil
`
	front = strings.ReplaceAll(front, "${sql}", sql)
	front = strings.ReplaceAll(front, "${table_name}", tableName)
	//填充model
	tmp := ""
	for _, sqlParam := range sqlParams {
		tmp += ",model." + sqlParam
	}
	front = strings.ReplaceAll(front, "${sql_params}", tmp[1:])
	front = strings.ReplaceAll(front, "${middle}", middleGenerate(src, hasGroup))
	return front
}

// middleGenerate 填充where
func middleGenerate(src interface{}, hasGroup bool) string {
	typeOf := reflect.TypeOf(src)
	tmp := ""
	for i := 0; i < typeOf.NumField(); i++ {
		if typeOf.Field(i).Name != "StartTime" && typeOf.Field(i).Name != "EndTime" && typeOf.Field(i).Name != "IsDownload" && typeOf.Field(i).Name != "LoginChannel" {
			tmp += fmt.Sprintf("\t\tWhere(\"%s=?\", model.%s).\n", addAndTransferToLower(typeOf.Field(i).Name), typeOf.Field(i).Name)
		}
		//单独对LoginChannel参数进行判断
		if typeOf.Field(i).Name == "LoginChannel" && !hasGroup {
			tmp += fmt.Sprintf("\t\tWhere(\"login_channel=?\", model.LoginChannel).\n")
		}
		if i == typeOf.NumField()-1 {
			tmp = tmp[:len(tmp)-1]
			if hasGroup {
				m := ""
				for j := 0; j < typeOf.NumField(); j++ {
					if typeOf.Field(j).Name != "StartTime" && typeOf.Field(j).Name != "EndTime" && typeOf.Field(j).Name != "IsDownload" && typeOf.Field(j).Name != "LoginChannel" {
						m += "," + addAndTransferToLower(typeOf.Field(j).Name)
					}
				}
				tmp += fmt.Sprintf(`Group("call_date%s").`, m)
			}
			tmp += fmt.Sprintf("Count(&count).Error; err != nil {")
		}
	}
	return tmp
}

// openAndWriteStringFile 将内容写入文件中
func openAndWriteStringFile(pathFile string, write string) error {

	fmt.Println("path:", pathFile)
	file, e := os.OpenFile(pathFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 6)
	defer file.Close()
	if e != nil {
		return e
	}
	if _, e = file.WriteString(write); e != nil {
		return e
	}
	return nil
}

// writeService 写入service
// ${name_dir} 包名
// ${req_struct} 请求结构体
// ${req_struct_name} 请求结构体名称
// ${struct_name} 结构体名称
// ${transfer_param} 需要传递的参数
// ${is_download} 是否有isdownload参数
func (r ReplaceType) writeService() error {
	var front = `package ${name_dir}

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "time"
    "fmt"
)



func ${struct_name}QueryList(c *gin.Context) {

    defer ${log_print}("${struct_name}QueryList", c, time.Now())
  	// 请求结构体
	${req_struct}
  
    req := ${req_struct_name}{}
  
    if err := c.ShouldBindJSON(&req); err != nil {
        logging.Errorf("${req_struct_name},err:%v\n", err.Error())
        c.JSON(http.StatusOK, gin.H{"error": -1, "message": err.Error()})
        return
    }
  
  // model 层函数 get + 当前方法名
  
    data, count, err := Get${struct_name}QueryList(${struct_name}ModuleReq{
${transfer_param}
	})
  
    if err != nil {
        logging.Errorf("Get${struct_name}List,err:%v\n", err.Error())
        c.JSON(http.StatusOK, gin.H{"error": -1, "message": err.Error()})
        return
    }
  	
	${is_download}
}
`
	front = strings.ReplaceAll(front, "${name_dir}", r.NameDir)
	front = strings.ReplaceAll(front, "${log_print}", r.LogPrint)
	front = strings.ReplaceAll(front, "${req_struct}", formatPrint(r.ReqSrc, r.NotesReqSrc))
	front = strings.ReplaceAll(front, "${req_struct_name}", r.NameReqSrc)
	front = strings.ReplaceAll(front, "${struct_name}", r.NameSrc)
	front = strings.ReplaceAll(front, "${transfer_param}", transferParam(r.ReqSrc))
	front = strings.ReplaceAll(front, "${is_download}", r.checkIsDownload(isParam(r.ReqSrc, "IsDownload")))

	if e := openAndWriteStringFile(path.Join(r.Path, r.NameDir, "service.go"), front); e != nil {
		return e
	}
	return nil
}

// checkIsDownload 判断是否有IsDownload
// ${notes} 数据库结构注释
// ${to_string} 将不同类型转换为string
// ${excel_table_name} 调用方法表名参数
func (r ReplaceType) checkIsDownload(hasIsd bool) string {
	if hasIsd {
		var front = `
if !req.IsDownload {
	c.JSON(http.StatusOK, gin.H{"error": 0, "message": "", "data": data, "count": count})
	return
} else {
	TableTitle := ${notes} 		// 数据库结构注释
	var TableDate [][]string	// 数据

	for _, v := range data {
		TableDate = append(TableDate, []string{
${to_string}
		}) //  根据数据库表结构 和数据的类型进行数据区分 转换成 string
	}

	//调用方法  表名参数
	utils.ExportToExcel(c, TableTitle, TableDate, "${excel_table_name}")

}
`
		front = strings.ReplaceAll(front, "${notes}", fmt.Sprintf("%#v", r.NotesSrc))
		front = strings.ReplaceAll(front, "${to_string}", strconvToString(r.Src))
		front = strings.ReplaceAll(front, "${excel_table_name}", r.ExcelTableName)
		return front
	} else {
		var front = `
	c.JSON(http.StatusOK, gin.H{"error": 0, "message": "", "data": data, "count": count})
	return
`
		return front
	}

}

// strconvToString 未完成
func strconvToString(src interface{}) string {
	typeOf := reflect.TypeOf(src)
	tmp := ""
	for i := 0; i < typeOf.NumField(); i++ {
		switch typeOf.Field(i).Type.String() {
		case "string":
			tmp += "\t\t\tv." + typeOf.Field(i).Name + ",\n"
		case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32":
			tmp += "\t\t\tstrconv.Itoa(int(v." + typeOf.Field(i).Name + ")),\n"
		case "int64":
			tmp += "\t\t\tstrconv.FormatInt(v." + typeOf.Field(i).Name + ",10),\n"
		case "uint64":
			tmp += "\t\t\tstrconv.FormatUint(v." + typeOf.Field(i).Name + ",10),\n"
		case "float32", "float64":
			tmp += "\t\t\tstrconv.FormatFloat(float64(v." + typeOf.Field(i).Name + "),'f',6,64),\n"
		case "bool":
			tmp += "\t\t\tstrconv.FormatBool(" + typeOf.Field(i).Name + "),\n"
		case "interface {}", "time.Time", "[]string", "[]int", "[]int32", "[]int64", "[]float32", "[]float64":
			tmp += `			fmt.Sprintf("%s",` + typeOf.Field(i).Name + "),\n"
		default:
		}
	}
	return tmp
}

// transferParam 将req的参数特定输出
func transferParam(src interface{}) string {
	typeOf := reflect.TypeOf(src)
	tmp := ""
	for i := 0; i < typeOf.NumField(); i++ {
		tmp += fmt.Sprintf("\t\t%s:    req.%s,\n", typeOf.Field(i).Name, typeOf.Field(i).Name)
	}
	return tmp
}

// writeRouter 路由编写
// ${name_dir} 包名
// ${struct_name} 结构体名称
// ${path_route} 路由路径
func (r ReplaceType) writeRouter() error {
	var front = `package ${name_dir}

import "github.com/gin-gonic/gin"

func Router(r gin.IRouter) {
	r.POST("${path_route}", ${struct_name}QueryList)   
}
`
	split := strings.Split(addAndTransferToLower(r.NameSrc), "_")
	tmp := "/"
	for _, v := range split {
		tmp += v + "/"
	}
	tmp += "query-list"

	front = strings.ReplaceAll(front, "${name_dir}", r.NameDir)
	front = strings.ReplaceAll(front, "${struct_name}", r.NameSrc)
	front = strings.ReplaceAll(front, "${path_route}", tmp)

	if e := openAndWriteStringFile(path.Join(r.Path, r.NameDir, "router.go"), front); e != nil {
		return errorx.Wrap(e)
	}
	return nil
}

// formatPrint 格式化输出结构体
// 问题：未结构化生成结构体
func formatPrint(src interface{}, notes []string) string {

	typeOf := reflect.TypeOf(src)
	tmp := ""
	for i := 0; i < typeOf.NumField(); i++ {
		if len(notes) > i {
			tmp += fmt.Sprintf("    %s  %s    `%s` // %s  \n", typeOf.Field(i).Name, typeOf.Field(i).Type, typeOf.Field(i).Tag, notes[i])
		} else {
			tmp += fmt.Sprintf("    %s  %s    `%s` \n", typeOf.Field(i).Name, typeOf.Field(i).Type, typeOf.Field(i).Tag)
		}
	}
	return fmt.Sprintf("type %s struct {\n%s \n}", typeOf.Name(), tmp)
}

// addAndTransferToLower
//现将第一个字符转换为小写，大写前添加下划线并全部转为小写,
func addAndTransferToLower(name string) string {
	name = strings.ToLower(string(name[0])) + name[1:]
	index := 0
	for k, v := range name {
		if v >= 65 && v <= 90 {
			name = name[0:k+index] + "_" + name[k+index:]
			index++
		}
	}
	return strings.ToLower(name)
}
