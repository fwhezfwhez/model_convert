package model_convert

import (
	"bufio"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strings"
)

const (
	JSON = 1 << iota
	FORM
	GORM
)

type ModelConvert struct {
	Db         *gorm.DB
	PrintModel bool
	TagFlag    int
}

// init modelConvert's dataSource
func (mc *ModelConvert) Init(dialect, dataSource string, printModel bool) error {
	var er error
	mc.Db, er = gorm.Open(dialect, dataSource)
	mc.PrintModel = printModel
	return er
}

// mc.SetFlag(model_convert.JSON | model_convert.GORM | model_convert.FORM)
func (mc *ModelConvert) SetFlags(flag int) {
	mc.TagFlag = flag
}
func (mc ModelConvert) HasTag() bool {
	ok1 := mc.TagFlag&JSON != 0
	ok2 := mc.TagFlag&GORM != 0
	ok3 := mc.TagFlag&FORM != 0
	if !ok1 && !ok2 && !ok3 {
		return false
	}
	return true
}

// close db
func (mc *ModelConvert) Destroy() {
	mc.Db.Close()
}

// point: 指定tag名，指定tag种类,若未指定表，则生成该数据库内所有表model
func (mc *ModelConvert) Generate(tables ...string) (string, error) {
	// 弱表未指定，则找到该数据库里所有表
	if tables == nil || len(tables) == 0 {
		tables = make([]string, 0, 5)
		tableSql := `
			SELECT
				-- c.relkind AS type,
				c.relname::varchar AS table_name
			FROM pg_class c
			JOIN ONLY pg_namespace n ON n.oid = c.relnamespace
			WHERE n.nspname = 'public'
			AND c.relkind = 'r'
			ORDER BY c.relname
		`
		mc.Db.SingularTable(true)
		if e := mc.Db.Raw(tableSql).Pluck("table_name", &tables).Error; e != nil {
			return "", errorx.New(e)
		}
		if len(tables) == 0 {
			return "", errorx.NewFromString("未在该dataSource下找到表，表数量为0")
		}
	}

	// 输出model并返回
	var FindColumnsSql = `
        SELECT
            a.attnum AS column_number,
            a.attname AS column_name,
            --format_type(a.atttypid, a.atttypmod) AS column_type,
            a.attnotnull AS not_null,
			COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') AS default_value,
    		COALESCE(ct.contype = 'p', false) AS  is_primary_key,
    		CASE
        	WHEN a.atttypid = ANY ('{int,int8,int2}'::regtype[])
          		AND EXISTS (
				SELECT 1 FROM pg_attrdef ad
             	WHERE  ad.adrelid = a.attrelid
             	AND    ad.adnum   = a.attnum
             	AND    ad.adsrc = 'nextval('''
                	|| (pg_get_serial_sequence (a.attrelid::regclass::text
                	                          , a.attname))::regclass
                	|| '''::regclass)'
             	)
            THEN CASE a.atttypid
                    WHEN 'int'::regtype  THEN 'serial'
                    WHEN 'int8'::regtype THEN 'bigserial'
                    WHEN 'int2'::regtype THEN 'smallserial'
                 END
			WHEN a.atttypid = ANY ('{uuid}'::regtype[]) AND COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') != ''
            THEN 'autogenuuid'
        	ELSE format_type(a.atttypid, a.atttypmod)
    		END AS column_type
		FROM pg_attribute a
		JOIN ONLY pg_class c ON c.oid = a.attrelid
		JOIN ONLY pg_namespace n ON n.oid = c.relnamespace
		LEFT JOIN pg_constraint ct ON ct.conrelid = c.oid
		AND a.attnum = ANY(ct.conkey) AND ct.contype = 'p'
		LEFT JOIN pg_attrdef ad ON ad.adrelid = c.oid AND ad.adnum = a.attnum
		WHERE a.attisdropped = false
		AND n.nspname = 'public'
		AND c.relname = ?
		AND a.attnum > 0
		ORDER BY a.attnum
	`
	var line, columnString, model, result string
	for _, tableName := range tables {
		var columns = make([]Column, 0, 10)
		if e := mc.Db.Raw(FindColumnsSql, tableName).Find(&columns).Error; e != nil {
			return "", errorx.New(e)
		}
		if ! mc.HasTag() {
			for _, column := range columns {
				line = fmt.Sprintf("    %s  %s\n", UnderLineToHump(column.ColumnName), typeConvert(column.ColumnType))
				columnString = columnString + line
			}
			model = fmt.Sprintf("type %s struct{\n%s}\n", UnderLineToHump(HumpToUnderLine(tableName)), columnString)
		} else {
			for _, column := range columns {
				line = fmt.Sprintf("    %s  %s    `", UnderLineToHump(column.ColumnName), typeConvert(column.ColumnType))
				if mc.TagFlag&GORM != 0 {
					line += fmt.Sprintf("gorm:\"column:%s\" ", column.ColumnName)
				}
				if mc.TagFlag&JSON != 0 {
					line += fmt.Sprintf("json:\"column:%s\" ", column.ColumnName)
				}
				if mc.TagFlag&FORM != 0 {
					line += fmt.Sprintf("form:\"column:%s\" ", column.ColumnName)
				}
				line += "`\n"
				columnString = columnString + line
			}

			model = fmt.Sprintf("type %s struct{\n%s}\n\n", UnderLineToHump(HumpToUnderLine(tableName)), columnString)

		}
		result += model
	}

	if mc.PrintModel {
		fmt.Println(result)
	}
	return result, nil
}

// AddJSONFormGormTag 添加json格式
func AddJSONFormGormTag(in string) string {
	var result string
	scanner := bufio.NewScanner(strings.NewReader(in))
	var oldLineTmp = ""
	var lineTmp = ""
	var propertyTmp = ""
	var seperateArr []string
	for scanner.Scan() {
		oldLineTmp = scanner.Text()
		lineTmp = strings.Trim(scanner.Text(), " ")
		if strings.Contains(lineTmp, "{") || strings.Contains(lineTmp, "}") {
			result = result + oldLineTmp + "\n"
			continue
		}
		seperateArr = Split(lineTmp, " ")
		// 接口或者父类声明不参与tag, 自带tag不参与tag
		if len(seperateArr) == 1 || len(seperateArr) == 3 {
			continue
		}
		propertyTmp = HumpToUnderLine(seperateArr[0])
		oldLineTmp = oldLineTmp + fmt.Sprintf("    `gorm:\"column:%s\" json:\"%s\" form:\"%s\"`", propertyTmp, propertyTmp, propertyTmp)
		result = result + oldLineTmp + "\n"
	}
	return result
}

// 根据数据源，表明获取列属性
func FindColumns(dataSource string, tableName string) []Column {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(fmt.Sprintf("recover from a fatal error : %v", e))
		}
	}()
	var FindColumnsSql = `
        SELECT
            a.attnum AS column_number,
            a.attname AS column_name,
            --format_type(a.atttypid, a.atttypmod) AS column_type,
            a.attnotnull AS not_null,
			COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') AS default_value,
    		COALESCE(ct.contype = 'p', false) AS  is_primary_key,
    		CASE
        	WHEN a.atttypid = ANY ('{int,int8,int2}'::regtype[])
          		AND EXISTS (
				SELECT 1 FROM pg_attrdef ad
             	WHERE  ad.adrelid = a.attrelid
             	AND    ad.adnum   = a.attnum
             	AND    ad.adsrc = 'nextval('''
                	|| (pg_get_serial_sequence (a.attrelid::regclass::text
                	                          , a.attname))::regclass
                	|| '''::regclass)'
             	)
            THEN CASE a.atttypid
                    WHEN 'int'::regtype  THEN 'serial'
                    WHEN 'int8'::regtype THEN 'bigserial'
                    WHEN 'int2'::regtype THEN 'smallserial'
                 END
			WHEN a.atttypid = ANY ('{uuid}'::regtype[]) AND COALESCE(pg_get_expr(ad.adbin, ad.adrelid), '') != ''
            THEN 'autogenuuid'
        	ELSE format_type(a.atttypid, a.atttypmod)
    		END AS column_type
		FROM pg_attribute a
		JOIN ONLY pg_class c ON c.oid = a.attrelid
		JOIN ONLY pg_namespace n ON n.oid = c.relnamespace
		LEFT JOIN pg_constraint ct ON ct.conrelid = c.oid
		AND a.attnum = ANY(ct.conkey) AND ct.contype = 'p'
		LEFT JOIN pg_attrdef ad ON ad.adrelid = c.oid AND ad.adnum = a.attnum
		WHERE a.attisdropped = false
		AND n.nspname = 'public'
		AND c.relname = ?
		AND a.attnum > 0
		ORDER BY a.attnum
	`
	db, err := gorm.Open("postgres", dataSource)
	db.SingularTable(true)
	//db.LogMode(true)
	if err != nil {
		panic(err)
	}
	var columns = make([]Column, 0, 10)
	db.Raw(FindColumnsSql, tableName).Find(&columns)
	return columns
}

// 数据库表转go model
func TableToStruct(dataSource string, tableName string) string {
	columnString := ""
	tmp := ""
	columns := FindColumns(dataSource, tableName)
	for _, column := range columns {

		tmp = fmt.Sprintf("    %s  %s\n", UnderLineToHump(column.ColumnName), typeConvert(column.ColumnType))
		columnString = columnString + tmp
	}

	rs := fmt.Sprintf("type %s struct{\n%s}", UnderLineToHump(HumpToUnderLine(tableName)), columnString)
	return rs
}

// 数据库表转go model 带tag
func TableToStructWithTag(dataSource string, tableName string) string {
	columnString := ""
	tmp := ""
	columns := FindColumns(dataSource, tableName)
	for _, column := range columns {

		tmp = fmt.Sprintf("    %s  %s    `gorm:\"column:%s;default:\" json:\"%s\" form:\"%s\"`\n",
			UnderLineToHump(column.ColumnName), typeConvert(column.ColumnType), column.ColumnName, column.ColumnName, column.ColumnName)
		columnString = columnString + tmp
	}
	var prefix = `
import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/errorx"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

/* 
    Code is auto-generated by github.com/fwhezfwhez/model_convert.Package below might be imported:
      - github.com/fwhezfwhez/errorx
      - github.com/garyburd/redigo/redis
      - github.com/jinzhu/gorm
    You can get them by:
      - go get github.com/fwhezfwhez/errorx
      - go get github.com/garyburd/redigo/redis
      - go get github.com/jinzhu/gorm

    To fulfill redis part, don't forget to set TODOs.They are:
      - RedisKey() string
      - RedisSecondDuration() int
*/
`

	var extra = `
func (o ${structName}) RedisKey() string {
	// TODO set its redis key
	return ""
}

func  (o ${structName}) RedisSecondDuration() int {
    // TODO set its redis duration, default one day, -1 without time limit
    return 1 * 24 * 60 * 60
}


func (o *${structName}) GetFromRedis(conn redis.Conn) error {
	if o.RedisKey() == "" {
		return errorx.NewFromString("object ${structName} has not set redis key yet")
	}
	buf,e:= redis.Bytes(conn.Do("GET", o.RedisKey()))

    if e==nil && string(buf)=="DISABLE"{
        return fmt.Errorf("not found record in db nor redis")
    }

	if e == redis.ErrNil {
		return e
	}

	if e != nil && e != redis.ErrNil {
		return errorx.Wrap(e)
	}

	e = json.Unmarshal(buf, &o)

	if e!=nil {
		return errorx.Wrap(e)
	}
	return nil
}

// engine should prepare its condition.
// if record not found,it will return 'var notFound = fmt.Errorf("not found record in db nor redis")'.
// If you want to ignore not found error, do it like:
// if e:= o.MustGet(conn, engine.Model(Model{}).Where("condition =?", arg)).Error;e!=nil {
//     if e.Error() == "not found record in db nor redis"{
//         log.Println(e)
//         return
//     }
// }
func (o *${structName}) MustGet(conn redis.Conn, engine *gorm.DB) error {
	e := o.GetFromRedis(conn)
    // When redis key stores its value 'DISABLE', will returns notFoundError and no need to query from db any more
    if e!=nil && e.Error() == "not found record in db nor redis" {
       return e
    }

	if e == nil {
		return nil
	}
	if e != nil {
		var count int
		if e2 := engine.Count(&count).Error; e2 != nil {
			return errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e2))
		}
		if count == 0 {
			var notFound = fmt.Errorf("not found record in db nor redis")
            conn.Do("SET", o.RedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			return notFound
		}

		if e3 := engine.First(&o).Error; e3 != nil {
			return errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e3))
		}
		if e == redis.ErrNil {
			o.SyncToRedis(conn)
			return nil
		}
		return errorx.Wrap(e)
	}
	return nil
}

func (o ${structName}) SyncToRedis(conn redis.Conn) error {
	if o.RedisKey() == "" {
		return errorx.NewFromString("object ${structName} has not set redis key yet")
	}
	buf, e := json.Marshal(o)
	if e != nil {
		return errorx.Wrap(e)
	}
	if o.RedisSecondDuration() == -1 {
		if _, e := conn.Do("SET", o.RedisKey(), buf); e != nil {
			return errorx.Wrap(e)
		}
	} else {
		if _, e := conn.Do("SETEX", o.RedisKey(), o.RedisSecondDuration(), buf); e != nil {
			return errorx.Wrap(e)
		}
	}
	return nil
}

func (o ${structName}) DeleteFromRedis(conn redis.Conn) error{
	if o.RedisKey() == "" {
		return errorx.NewFromString("object ${structName} has not set redis key yet")
	}
	if _, e := conn.Do("DEL", o.RedisKey()); e != nil {
		return errorx.Wrap(e)
	}
	return nil
}
`
	// extra = strings.ReplaceAll(extra,"${structName}", UnderLineToHump(HumpToUnderLine(tableName)))
	extra = strings.Replace(extra, "${structName}", UnderLineToHump(HumpToUnderLine(tableName)), -1)
	rs := fmt.Sprintf("%stype %s struct{\n%s}\n\nfunc (o %s) TableName() string {\n    return \"%s\" \n}\n\n%s", prefix, UnderLineToHump(HumpToUnderLine(tableName)), columnString, UnderLineToHump(tableName), tableName, extra)
	return rs
}

// UnderLineToHump 下划线转驼峰
func UnderLineToHump(s string) string {
	arr := strings.Split(s, "_")
	for i, v := range arr {
		arr[i] = strings.ToUpper(string(v[0])) + string(v[1:])
	}
	return strings.Join(arr, "")
}

// 类型转换pg->go
func typeConvert(s string) string {
	if strings.Contains(s, "char") || in(s, []string{
		"text",
	}) {
		return "string"
	}
	if in(s, []string{"bigint", "bigserial", "integer", "smallint", "serial", "big serial"}) {
		return "int"
	}
	if in(s, []string{"numeric", "decimal", "real"}) {
		return "decimal.Decimal"
	}
	if in(s, []string{"bytea"}) {
		return "[]byte"
	}
	if strings.Contains(s, "time") || in(s, []string{"date"}) {
		return "time.Time"
	}
	if in(s, []string{"jsonb"}) {
		return "json.RawMessage"
	}
	return "interface{}"
}

// s 是否in arr
func in(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

// Split 增强型Split，对  a,,,,,,,b,,c     以","进行切割成[a,b,c]
func Split(s string, sub string) []string {
	var rs = make([]string, 0, 20)
	tmp := ""
	Split2(s, sub, &tmp, &rs)
	return rs
}

// Split2 附属于Split，可独立使用
func Split2(s string, sub string, tmp *string, rs *[]string) {
	s = strings.Trim(s, sub)
	if !strings.Contains(s, sub) {
		*tmp = s
		*rs = append(*rs, *tmp)
		return
	}
	for i := range s {
		if string(s[i]) == sub {
			*tmp = s[:i]
			*rs = append(*rs, *tmp)
			s = s[i+1:]
			Split2(s, sub, tmp, rs)
			return
		}
	}
}

// HumpToUnderLine 驼峰转下划线
func HumpToUnderLine(s string) string {
	if s == "ID" {
		return "id"
	}
	var rs string
	elements := FindUpperElement(s)
	for _, e := range elements {
		s = strings.Replace(s, e, "_"+strings.ToLower(e), -1)
	}
	rs = strings.Trim(s, " ")
	rs = strings.Trim(rs, "\t")
	return strings.Trim(rs, "_")
}

// FindUpperElement 找到字符串中大写字母的列表,附属于HumpToUnderLine
func FindUpperElement(s string) []string {
	var rs = make([]string, 0, 10)
	for i := range s {
		if s[i] >= 65 && s[i] <= 90 {
			rs = append(rs, string(s[i]))
		}
	}
	return rs
}

// 数据库列属性
type Column struct {
	ColumnNumber int    `gorm:"column_number"` // column index
	ColumnName   string `gorm:"column_name"`   // column_name
	ColumnType   string `gorm:"column_type"`   // column_type
}

// go model 带tag
func AddJSONFormTag(s string) string {
	var result string
	scanner := bufio.NewScanner(strings.NewReader(s))
	var oldLineTmp = ""
	var lineTmp = ""
	var propertyTmp = ""
	var seperateArr []string
	for scanner.Scan() {
		oldLineTmp = scanner.Text()
		lineTmp = strings.Trim(scanner.Text(), " ")
		if strings.Contains(lineTmp, "{") || strings.Contains(lineTmp, "}") {
			result = result + oldLineTmp + "\n"
			continue
		}
		seperateArr = Split(lineTmp, " ")
		// 接口或者父类声明不参与tag, 自带tag不参与tag
		if len(seperateArr) == 1 || len(seperateArr) == 3 {
			continue
		}
		propertyTmp = HumpToUnderLine(seperateArr[0])
		oldLineTmp = oldLineTmp + fmt.Sprintf("    `json:\"%s\" form:\"%s\"`", propertyTmp, propertyTmp)
		result = result + oldLineTmp + "\n"
	}
	return result
}
