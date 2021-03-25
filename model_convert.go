package model_convert

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/fwhezfwhez/errorx"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
		if !mc.HasTag() {
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
func FindColumns(dialect string, dataSource string, tableName string) []Column {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(fmt.Sprintf("recover from a fatal error : %v", e))
		}
	}()

	switch dialect {
	case "postgres", "pg", "psql":
		return findPGColumns(dataSource, tableName)
	case "mysql":
		return findMysqlColumns(dataSource, tableName)
	}
	return nil
}
func findPGColumns(dataSource string, tableName string) []Column {
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
             	-- AND    ad.adsrc = 'nextval('''
                --	|| (pg_get_serial_sequence (a.attrelid::regclass::text
                --	                          , a.attname))::regclass
                --	|| '''::regclass)'
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
func findMysqlColumns(dataSource string, tableName string) []Column {
	var FindColumnsSql = `
       SELECT column_name as column_name, column_type as column_type  FROM information_schema.columns WHERE table_name= ?
	`
	db, err := gorm.Open("mysql", dataSource)
	db.SingularTable(true)
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	var columns = make([]Column, 0, 10)
	if e := db.Raw(FindColumnsSql, tableName).Scan(&columns).Error; e != nil {
		panic(err)
	}

	return columns
}

// 数据库表转go model
func TableToStruct(dataSource string, tableName string, dialect ...string) string {
	if len(dialect) == 0 {
		dialect = []string{"postgres"}
	}
	columnString := ""
	tmp := ""
	columns := FindColumns(dialect[0], dataSource, tableName)
	for _, column := range columns {

		tmp = fmt.Sprintf("    %s  %s\n", UnderLineToHump(column.ColumnName), typeConvert(column.ColumnType))
		columnString = columnString + tmp
	}

	rs := fmt.Sprintf("type %s struct{\n%s}", UnderLineToHump(HumpToUnderLine(tableName)), columnString)
	return rs
}

type TableToStructWithTagReplacement struct {
	DBInstance    string
	DBInstancePkg string
}

func (o *TableToStructWithTagReplacement) Init() {
	o.DBInstance = "db.DB"
	o.DBInstancePkg = "path/to/db"
}

// 数据库表转go model 带tag
// replacements 取值类型为 string 或者 map[string]interface{}
// 取值为string时，replacements[0]为数据库类型,取值如postgres, mysql。并且，"postgres" 等价于map[string]interface{}{"dialect":"postgres"}
// 取值map时，将可以附加一些参数
func TableToStructWithTag(dataSource string, tableName string, replacements ...interface{}) string {

	columnString := ""
	tmp := ""

	dialect := getDialect(replacements...)

	ttswtr := fillTableToStructWithTagReplacement(replacements...)

	columns := FindColumns(dialect, dataSource, tableName)
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
     "${db_instance_pkg}"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

// Auto-Generate Header
/* 
    Code is auto-generated by github.com/fwhezfwhez/model_convert.Package below might be imported:
      - github.com/fwhezfwhez/errorx
      - github.com/garyburd/redigo/redis
      - github.com/jinzhu/gorm
      - ${db_instance_pkg}
    You can get them by:
      - go get github.com/fwhezfwhez/errorx
      - go get github.com/garyburd/redigo/redis
      - go get github.com/jinzhu/gorm

    To fulfill redis part, don't forget to set TODOs.They are:
      - RedisKey() string
      - RedisSecondDuration() int
*/
`
	prefix = strings.Replace(prefix, "${db_instance_pkg}", ttswtr.DBInstancePkg, -1)

	var extra = `
func (o ${structName}) DB() *gorm.DB {
    return ${db_instance}
}

var ${structName}RedisKeyFormat = ""

func (o ${structName}) RedisKey() string {
	// TODO set its redis key and required args
	return fmt.Sprintf(${structName}RedisKeyFormat, )
}


var Array${structName}RedisKeyFormat = ""

func (o ${structName}) ArrayRedisKey() string {
	// TODO set its array key and required args
	return fmt.Sprintf(Array${structName}RedisKeyFormat,)
}

func  (o ${structName}) RedisSecondDuration() int {
    // TODO set its redis duration, default 1-7 day,  return -1 means no time limit
    return int(time.Now().Unix() % 7 + 1) * 24 * 60 * 60
}

// TODO,set using db or not. If set false, o.MustGet() will never get its data from db.
func (o ${structName}) UseDB() bool {
	return false
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

func (o *${structName}) ArrayGetFromRedis(conn redis.Conn) ([]${structName}, error) {
	if o.ArrayRedisKey() == "" {
		return nil, errorx.NewFromString("object ${structName} has not set redis key yet")
	}

	var list = make([]${structName}, 0, 10)
	buf, e := redis.Bytes(conn.Do("GET", o.ArrayRedisKey()))

	// avoid passing through and hit database
	// When o.ArrayMustGet() not found both in redis and db, will set its key DISABLE
	// and return 'fmt.Errorf("not found record in db nor redis")'
	if e == nil && string(buf) == "DISABLE" {
		return nil, fmt.Errorf("not found record in db nor redis")
	}

	// Not found in redis
	if e == redis.ErrNil {
		return nil, e
	}

	// Server error, should be logged by caller
	if e != nil && e != redis.ErrNil {
		return nil, errorx.Wrap(e)
	}

	e = json.Unmarshal(buf, &list)

	if e != nil {
		return nil, errorx.Wrap(e)
	}
	return list, nil
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
	var shouldSyncToCache bool

	if ${structName}CacheSwitch {
		if e := o.GetFromCache(); e == nil {
			return nil
		}
		defer func() {
			if shouldSyncToCache {
				fmt.Println("exec sync to cache")
				o.SyncToCache()
			}
		}()
	}

	e := o.GetFromRedis(conn)
    // When redis key stores its value 'DISABLE', will returns notFoundError and no need to query from db any more
    if e!=nil && e.Error() == "not found record in db nor redis" {
       return e
    }

	if e == nil {
		shouldSyncToCache = true
		return nil
	}
	if e != nil {
		var count int
		if e2 := engine.Count(&count).Error; e2 != nil {
			return errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e2))
		}
		if count == 0 {
			var notFound = fmt.Errorf("not found record in db nor redis")
			if o.RedisSecondDuration() == -1 {
				conn.Do("SET", o.RedisKey(), "DISABLE", "NX")
			} else {
				conn.Do("SET", o.RedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			}
	        return notFound
		}

		if e3 := engine.First(&o).Error; e3 != nil {
			return errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e3))
		}
		shouldSyncToCache = true

		if e == redis.ErrNil {
			o.SyncToRedis(conn)
			return nil
		}
		return errorx.Wrap(e)
	}
	return nil
}

func (o *${structName}) ArrayMustGet(conn redis.Conn, engine *gorm.DB) ([]${structName}, error) {
	var shouldSyncToCache bool
	var arr []${structName}

	if Array${structName}CacheSwitch {
		if arr, e := o.ArrayGetFromCache(); e == nil {
			return  arr, nil
		}
		defer func() {
			if shouldSyncToCache {
				fmt.Println("exec sync to cache")
				o.ArraySyncToCache(arr)
			}
		}()
	}


	list, e := o.ArrayGetFromRedis(conn)
	// When redis key stores its value 'DISABLE', will returns notFoundError and no need to query from db any more
	// When call ArrayDeleteFromRedis(), will activate its redis and db query
	if e != nil && e.Error() == "not found record in db nor redis" {
		return nil, e
	}
	// get from redis success.
	if e == nil {
		shouldSyncToCache = true
		arr = list
		return list, nil
	}
	// get from redis fail, try db
	if e != nil {
		var count int
		if e2 := engine.Count(&count).Error; e2 != nil {
			return nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e2))
		}
		if count == 0 {
			var notFound = fmt.Errorf("not found record in db nor redis")
			if o.RedisSecondDuration() == -1 {
				conn.Do("SET", o.ArrayRedisKey(), "DISABLE", "NX")
			} else {
				conn.Do("SET", o.ArrayRedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			}
			return nil, notFound
		}

		if e3 := engine.Find(&list).Error; e3 != nil {
			return nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e3))
		}

		shouldSyncToCache = true
		arr = list
		// try sync to redis
		if e == redis.ErrNil {
			o.ArraySyncToRedis(conn, list)
			return list, nil
		}
		return nil, errorx.Wrap(e)
	}
	return nil, nil
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

func (o ${structName}) ArraySyncToRedis(conn redis.Conn, list []${structName}) error {
	if o.ArrayRedisKey() == "" {
		return errorx.NewFromString("object ${structName} has not set redis key yet")
	}
	buf, e := json.Marshal(list)
	if e != nil {
		return errorx.Wrap(e)
	}
	if o.RedisSecondDuration() == -1 {
		if _, e := conn.Do("SET", o.ArrayRedisKey(), buf); e != nil {
			return errorx.Wrap(e)
		}
	} else {
		if _, e := conn.Do("SETEX", o.ArrayRedisKey(), o.RedisSecondDuration(), buf); e != nil {
			return errorx.Wrap(e)
		}
	}
	return nil
}

func (o ${structName}) DeleteFromRedis(conn redis.Conn) error{
	if o.RedisKey() != "" {
		if _, e := conn.Do("DEL", o.RedisKey()); e != nil {
		    return errorx.Wrap(e)
	    }
	}
    
    if o.ArrayRedisKey() != "" {
	    if _, e := conn.Do("DEL", o.ArrayRedisKey()); e != nil {
	    	return errorx.Wrap(e)
    	}
    }

	if ${structName}CacheSwitch {
		o.deleteFromCache()
	}
	if Array${structName}CacheSwitch {
		o.ArraydeleteFromCache()
	}
	return nil
}
func (o ${structName}) ArrayDeleteFromRedis(conn redis.Conn) error{
	return o.DeleteFromRedis(conn)
}

// Dump data through api GET remote url generated by 'GenerateListApi()' to local database.
// This method should never used in production. It's best to to run it before app is running.
//
// mode=1, each time will delete old local data and dump from api.
// mode=2, each time will update/keep the existed data. Mode=2 is developing.
func (o ${structName}) DumpToLocal(url string,engine *gorm.DB, mode int) error {
    tableName := o.TableName()

    tran := engine.Begin()
    if e:=tran.Exec(fmt.Sprintf("delete from %s", tableName)).Error; e!=nil{
        tran.Rollback()
        return errorx.Wrap(e)
    }

    type Result struct{
        Data []${structName} ${data_json_tag}
        Count int ${count_json_tag}
    }
    var result Result
    resp,e:=http.Get(url)
    if e!=nil {
        tran.Rollback()
        return errorx.Wrap(e)
    }
    if resp ==nil || resp.Body ==nil {
        tran.Rollback()
        return errorx.NewFromString("resp or body nil")
    }
    defer resp.Body.Close()

    buf,e := ioutil.ReadAll(resp.Body)
    if e!=nil {
        tran.Rollback()
        return errorx.Wrap(e)
    }

	if resp.StatusCode != 200 {
        var body string
        if len(buf)<100 {
            body = string(buf)
        } else{
            body = string(buf[:100])
        }
		return errorx.NewFromStringf("status not 200, got %d,body %s", resp.StatusCode, body)
	}

    if e:=json.Unmarshal(buf, &result);e!=nil{
        tran.Rollback()
        return errorx.Wrap(e)
    }

    for i,_ := range result.Data {
        data := result.Data[i]
        if e:=tran.Model(&o).Create(&data).Error; e!=nil{
            tran.Rollback()
            return errorx.Wrap(e)
        }
    }
    tran.Commit()
    return nil
}

// 2nd-cache Header
// 2nd-cache share RedisKey() as its key.

// self Header
var (
	${structName}Cache         = make(map[string]${structName}, 0)
	${structName}CacheKeyOrder = make([]string, 0, 500)

	${structName}CacheLock = sync.RWMutex{}

	${structName}NotFoundErr  = fmt.Errorf("not found in cache")
	${structName}SwitchOffErr = fmt.Errorf("2nd-cache switch is off")
)

const (
	// Max cap of map and len
	${structName}CacheMaxLength = 5000
	// When faced with max cap, will delete 2000/5 data from map and arr
	// DeleteRate can't be zero.
	${structName}DeleteRate = 5
	// Whether use cache
	${structName}CacheSwitch = false
)

func (o *${structName}) GetFromCache() error {
	if ${structName}CacheSwitch == false {
		return ${structName}SwitchOffErr
	}
	${structName}CacheLock.RLock()
	defer ${structName}CacheLock.RUnlock()
	tmp, ok := ${structName}Cache[o.RedisKey()]
	if !ok {
		return ${structName}NotFoundErr
	}
	*o = tmp
	fmt.Println("get from cache")
	return nil
}

func (o *${structName}) deleteFromCache() {
	if ${structName}CacheSwitch == false {
		return
	}
	${structName}CacheLock.Lock()
	defer ${structName}CacheLock.Unlock()

	delete(${structName}Cache, o.RedisKey())
}

func (o *${structName}) SyncToCache() {
	if ${structName}CacheSwitch == false {
		return
	}

	if ${structName}DeleteRate == 0 || ${structName}DeleteRate < 0 {
		return
	}

	if ${structName}CacheMaxLength == 0 {
		return
	}

	${structName}CacheLock.Lock()
	defer ${structName}CacheLock.Unlock()

	leng := len(${structName}CacheKeyOrder)
	if leng >= ${structName}CacheMaxLength {
		delta := ${structName}CacheMaxLength / ${structName}DeleteRate
		for i := 0; i < delta; i++ {
			if _, ok := ${structName}Cache[${structName}CacheKeyOrder[i]]; ok {
				delete(${structName}Cache, ${structName}CacheKeyOrder[i])
            }
		}
		${structName}CacheKeyOrder = ${structName}CacheKeyOrder[delta:]
	}
	${structName}Cache[o.RedisKey()] = *o
	${structName}CacheKeyOrder = append(${structName}CacheKeyOrder, o.RedisKey())
}

// self Tail

// array Header
var (
	Array${structName}Cache         = make(map[string][]${structName}, 0)
	Array${structName}CacheKeyOrder = make([]string, 0, 500)

	Array${structName}CacheLock = sync.RWMutex{}

	Array${structName}NotFoundErr  = fmt.Errorf("not found in cache")
	Array${structName}SwitchOffErr = fmt.Errorf("2nd-cache switch is off")
)

const (
	// Max cap of map and len
	Array${structName}CacheMaxLength = 5000
	// When faced with max cap, will delete 2000/5 data from map and arr
	// DeleteRate can't be zero.
	Array${structName}DeleteRate = 5
	// Whether use cache
	Array${structName}CacheSwitch = false
)

func (o *${structName}) ArrayGetFromCache() ([]${structName}, error) {
	if Array${structName}CacheSwitch == false {
		return nil, Array${structName}SwitchOffErr
	}
	Array${structName}CacheLock.RLock()
	defer Array${structName}CacheLock.RUnlock()
	tmp, ok := Array${structName}Cache[o.ArrayRedisKey()]
	if !ok {
		return nil, Array${structName}NotFoundErr
	}
	fmt.Println("get from cache")
	return tmp, nil
}

func (o *${structName}) ArraydeleteFromCache() {
	if Array${structName}CacheSwitch == false {
		return
	}
	Array${structName}CacheLock.Lock()
	defer Array${structName}CacheLock.Unlock()

	delete(Array${structName}Cache, o.ArrayRedisKey())
}

func (o *${structName}) ArraySyncToCache(arr []${structName}) {
	if Array${structName}CacheSwitch == false {
		return
	}

	if Array${structName}DeleteRate == 0 || Array${structName}DeleteRate < 0 {
		return
	}

	if Array${structName}CacheMaxLength == 0 {
		return
	}

	Array${structName}CacheLock.Lock()
	defer Array${structName}CacheLock.Unlock()

	leng := len(Array${structName}CacheKeyOrder)
	if leng >= Array${structName}CacheMaxLength {
		delta := Array${structName}CacheMaxLength / Array${structName}DeleteRate
		for i := 0; i < delta; i++ {
			if _, ok := Array${structName}Cache[Array${structName}CacheKeyOrder[i]]; ok {
				delete(Array${structName}Cache, Array${structName}CacheKeyOrder[i])
            }
		}
		Array${structName}CacheKeyOrder = Array${structName}CacheKeyOrder[delta:]
	}
	Array${structName}Cache[o.ArrayRedisKey()] = arr
	Array${structName}CacheKeyOrder = append(Array${structName}CacheKeyOrder, o.ArrayRedisKey())
}
// array Tail

// 2nd-cache Tail

// flexible-cache Header
// func (o ${structName}) ${cache_name}Key() string{
// 	// TODO-Set cache redis key
// 	return ""
// }
// func (o ${structName}) ${cache_name}Duration() int{
// 	// TODO-Set cache redis key expire duration. Default 1-7 days
//     return int(time.Now().Unix() % 7 + 1) * 24 * 60 * 60
// }
// func (o *${structName}) ${cache_name}MustGet(conn redis.Conn, source func(${cache_name} *${cache_type})error) (${cache_type}, error) {

// 	rs, e:= redis.${Cache_type}(conn.Do("GET", o.${cache_name}Key()))
// 	if e !=nil {
// 		if e == redis.ErrNil {
//             if e:=source(&rs); e!=nil {
// 				return rs, errorx.Wrap(e)
// 			}
// 			if _, e= conn.Do("SETEX",  o.${cache_name}Key(), ${cache_name}Duration(), rs),; e!=nil {
// 				return rs, errorx.Wrap(e)
// 			}
// 			return rs,nil
// 		}
// 		return rs, errorx.Wrap(e)
// 	}
// 	return rs,nil

// }
// flexible-cache Tail

// no-decode Header
// 
// MustGetNoDecode do most similar work as MustGet do, but it will not unmarshal data from redis into 'o', in the meanwhile, will return its raw json stream as return.
// This function aims to save cost of decoding in the only case that you want to return 'o' itself and has nothing changed to inner values.
// 'engine' should prepare its condition.
// if record not found,it will return 'var notFound = fmt.Errorf("not found record in db nor redis")'.
// If you want to ignore not found error, do it like:
// if buf, e:= o.MustGetNoDecode(conn, engine.Model(Model{}).Where("condition =?", arg)).Error;e!=nil {
//     if e.Error() == "not found record in db nor redis" || e == redis.ErrNil {
//         log.Println(e)
//         return
//     }
// }
// 
func (o *${structName}) MustGetNoDecode(conn redis.Conn, engine *gorm.DB) (json.RawMessage, error) {
	var shouldSyncToCache bool

	if ${structName}CacheSwitch {
		if e := o.GetFromCache(); e == nil {
			return nil, nil
		}
		defer func() {
			if shouldSyncToCache {
				fmt.Println("exec sync to cache")
				o.SyncToCache()
			}
		}()
	}

	arrBuf, e := o.GetFromRedisNoDecode(conn)
    // When redis key stores its value 'DISABLE', will returns notFoundError and no need to query from db any more
    if e!=nil && e.Error() == "not found record in db nor redis" {
       return nil, e
    }

	if e == nil {
		shouldSyncToCache = true
		return arrBuf, nil
	}
	if e != nil {
		var count int
		if e2 := engine.Count(&count).Error; e2 != nil {
			return nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e2))
		}
		if count == 0 {
			var notFound = fmt.Errorf("not found record in db nor redis")
			if o.RedisSecondDuration() == -1 {
				conn.Do("SET", o.RedisKey(), "DISABLE", "NX")
			} else {
				conn.Do("SET", o.RedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			}
	        return nil, notFound
		}

		if e3 := engine.First(&o).Error; e3 != nil {
			return nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e3))
		}
		shouldSyncToCache = true

		if e == redis.ErrNil {
			o.SyncToRedis(conn)
			return nil,nil 
		}
		return nil, errorx.Wrap(e)
	}
	return nil,nil
}

// GetFromRedisNoDecode will return its json raw stream and will not decode into 'o'.
// It aims to save cost of decoding if json stream is decoded slowly.
func (o *${structName}) GetFromRedisNoDecode(conn redis.Conn) (json.RawMessage, error) {
	if o.RedisKey() == "" {
		return nil, errorx.NewFromString("object ${structName} has not set redis key yet")
	}
	buf,e:= redis.Bytes(conn.Do("GET", o.RedisKey()))

    if e==nil && string(buf)=="DISABLE"{
        return nil, fmt.Errorf("not found record in db nor redis")
    }

	if e == redis.ErrNil {
		return nil, e
	}

	if e != nil && e != redis.ErrNil {
		return nil, errorx.Wrap(e)
	}

	return buf, nil
}

// ArrayMustGetNoDecode will not unmarshal json stream to 'arr' and return json.Rawmessage as return value instead if it's found in redis,
// otherwise will return arr from cache or db.
//
// This function aims to save cost of decoding in the read-only case of 'o'. It means you should do nothing changed to its json value. 
/* 
	arr, arrBuf, e:= o.ArrayMustGetNoDecode(conn, engine)
	if e!=nil {
	// handle error
	}

	if len(arrBuf) >0 {
	c.JSON(200, gin.H{"message":"success", "data": arrBuf})
	} else {
		c.JSON(200, gin.H{"message":"success", "data": arr})
	}
*/
func (o *${structName}) ArrayMustGetNoDecode(conn redis.Conn, engine *gorm.DB) ([]${structName},json.RawMessage, error) {
	var shouldSyncToCache bool
	var arr []${structName}

	if Array${structName}CacheSwitch {
		if arr, e := o.ArrayGetFromCache(); e == nil {
			return  arr, nil,  nil
		}
		defer func() {
			if shouldSyncToCache {
				fmt.Println("exec sync to cache")
				o.ArraySyncToCache(arr)
			}
		}()
	}


	arrBuf, e := o.ArrayGetFromRedisNoDecode(conn)
	// When redis key stores its value 'DISABLE', will returns notFoundError and no need to query from db any more
	// When call ArrayDeleteFromRedis(), will activate its redis and db query
	if e != nil && e.Error() == "not found record in db nor redis" {
		return nil,nil, e
	}
	// get from redis success.
	if e == nil {
		// shouldSyncToCache = true
		// arr = list
		return nil, arrBuf, nil
	}
	// get from redis fail, try db
	if e != nil {
		var list = make([]${structName},0, 100)
		var count int
		if e2 := engine.Count(&count).Error; e2 != nil {
			return nil,nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e2))
		}
		if count == 0 {
			var notFound = fmt.Errorf("not found record in db nor redis")
			if o.RedisSecondDuration() == -1 {
				conn.Do("SET", o.ArrayRedisKey(), "DISABLE", "NX")
			} else {
				conn.Do("SET", o.ArrayRedisKey(), "DISABLE", "EX", o.RedisSecondDuration(), "NX")
			}
			return nil, nil, notFound
		}

		if e3 := engine.Find(&list).Error; e3 != nil {
			return nil, nil, errorx.GroupErrors(errorx.Wrap(e), errorx.Wrap(e3))
		}

		shouldSyncToCache = true
		arr = list
		// try sync to redis
		if e == redis.ErrNil {
			o.ArraySyncToRedis(conn, list)
			return list,  nil, nil
		}
		return nil, nil,  errorx.Wrap(e)
	}
	return nil, nil, nil
}

func (o *${structName}) ArrayGetFromRedisNoDecode(conn redis.Conn) (json.RawMessage, error) {
	if o.ArrayRedisKey() == "" {
		return nil, errorx.NewFromString("object ${structName} has not set redis key yet")
	}

	buf, e := redis.Bytes(conn.Do("GET", o.ArrayRedisKey()))

	// avoid passing through and hit database
	// When o.ArrayMustGet() not found both in redis and db, will set its key DISABLE
	// and return 'fmt.Errorf("not found record in db nor redis")'
	if e == nil && string(buf) == "DISABLE" {
		return nil, fmt.Errorf("not found record in db nor redis")
	}

	// Not found in redis
	if e == redis.ErrNil {
		return nil, e
	}

	// Server error, should be logged by caller
	if e != nil && e != redis.ErrNil {
		return nil, errorx.Wrap(e)
	}

	return buf, nil
}
// no-decode Tail
// Auto-Generate Tail
`

	// extra = strings.ReplaceAll(extra,"${structName}", UnderLineToHump(HumpToUnderLine(tableName)))
	extra = strings.Replace(extra, "${db_instance}", ttswtr.DBInstance, -1)
	extra = strings.Replace(extra, "${db_instance_pkg}", ttswtr.DBInstancePkg, -1)

	extra = strings.Replace(extra, "${structName}", UnderLineToHump(HumpToUnderLine(tableName)), -1)
	extra = strings.Replace(extra, "${count_json_tag}", "`json:\"count\"`", -1)
	extra = strings.Replace(extra, "${data_json_tag}", "`json:\"data\"`", -1)

	rs := fmt.Sprintf("%stype %s struct{\n%s}\n\nfunc (o %s) TableName() string {\n    return \"%s\" \n}\n\n%s", prefix, UnderLineToHump(HumpToUnderLine(tableName)), columnString, UnderLineToHump(tableName), tableName, extra)

	// 增加注入剪贴板
	clipboard.WriteAll(rs)
	return rs
}

func getDialect(replacements ...interface{}) string {
	if len(replacements) == 0 {
		return "postgres"
	}

	var replacementI interface{}
	if len(replacements) == 1 {
		dialect, convert := replacements[0].(string)
		if convert {
			return dialect
		}
	}

	replacementI = replacements[0]

	switch v := replacementI.(type) {
	case map[string]interface{}:
		dialect, exist := v["dialect"]
		if exist {
			dialectStr, assert := dialect.(string)
			if assert {
				return dialectStr
			}
			panic(fmt.Errorf("replacements[0]['dialect'] must be string type,but got %v", dialect))
		} else {
			return "postgres"
		}
	default:
		panic(fmt.Errorf("replacements[0] only accepts string or map[string]interface{}"))
	}
}

func fillTableToStructWithTagReplacement(replacements ...interface{}) TableToStructWithTagReplacement {
	var ttstr TableToStructWithTagReplacement
	ttstr.Init()

	if len(replacements) == 0 {
		return ttstr
	}
	replacement := replacements[0]

	switch v := replacement.(type) {
	case map[string]interface{}:
		di, exist := v["${db_instance}"]
		if exist {
			ttstr.DBInstance = fmt.Sprintf("%v", di)
		}
		dip, exist := v["${db_instance_pkg}"]
		if exist {
			ttstr.DBInstancePkg = fmt.Sprintf("%v", dip)
		}
		return ttstr
	case string:
		return ttstr
	default:
		return ttstr
	}
}

// UnderLineToHump 下划线转驼峰
func UnderLineToHump(s string) string {
	arr := strings.Split(s, "_")
	for i, v := range arr {
		arr[i] = strings.ToUpper(string(v[0])) + string(v[1:])
	}
	return strings.Join(arr, "")
}

func LowerFirstLetter(s string) string {
	return strings.ToLower(string(s[0])) + string(s[1:])
}

// 类型转换pg->go
func typeConvert(s string) string {
	if strings.Contains(s, "char") || in(s, []string{
		"text",
	}) {
		return "string"
	}
	// postgres
	{
		if in(s, []string{"double precision", "double"}) {
			return "float64"
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
		if strings.Contains(s, "time") || in(s, []string{"date", "datetime", "timestamp"}) {
			return "time.Time"
		}
		if in(s, []string{"jsonb"}) {
			return "json.RawMessage"
		}
		if in(s, []string{"bool", "boolean"}) {
			return "bool"
		}

		if in(s, []string{"bigint[]"}) {
			return "[]int64"
		}
	}
	// mysql
	{
		if strings.HasPrefix(s, "int") {
			return "int"
		}
		if strings.HasPrefix(s, "varchar") {
			return "string"
		}
		if s == "json" {
			return "json.RawMessage"
		}
		if in(s, []string{"bool", "boolean"}) {
			return "bool"
		}
	}

	return s
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
	ColumnName string `gorm:"column:column_name"` // column_name
	ColumnType string `gorm:"column:column_type"` // column_type
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
