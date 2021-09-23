package model_convert

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/atotto/clipboard"
)

// GenerateListWhere requires a model and generate api codes to list the array of the model.
// It can fit dynamic where occasion, page, size, order-by.
//
// Only support gin-gorm.
//
// src is a go model instance.
//
// withListArgs will generate page,size, order by and its db engine count,find.
//
// Replacement optional as:
// - ${jump_fields}, "password,pw"
// - ${layout}, "2006-01-02 15:04:03"
// - ${time_zone} "time.Local"
func GenerateListWhere(src interface{}, withListArgs bool, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	var format = `
   ${handle}
`
	vValue := reflect.ValueOf(src)
	vType := reflect.TypeOf(src)
	handle := func() string {
		var result string
		for i := 0; i < vValue.NumField(); i++ {
			// only handle basic types, otherwise continue
			if !in(vType.Field(i).Type.String(), []string{
				"int", "int8", "int16", "int32", "int64",
				"float32", "float64",
				"string",
				"uint8", "uint16", "uint32", "uint64",
				"time.Time",
			}) {
				continue
			}

			// handle time.Time alone
			if vType.Field(i).Type.AssignableTo(reflect.TypeOf(time.Time{})) {
				tagName := HumpToUnderLine(vType.Field(i).Name)
				fieldName := LowerFistLetter(vType.Field(i).Name)
				zeroValue := GetZeroValue(vValue.Field(i).Interface())
				layout := replacement[0]["layout"]
				if layout == "" {
					layout = "2006-01-02 15:04:05"
				}
				listLayout := replacement[0]["${list_layout}"]
				if layout == "" {
					layout = "2006-01-02"
				}
				timeZone := replacement[0]["time_zone"]
				if timeZone == "" {
					timeZone = "time.Local"
				}
				tmp := `
    ${tag_name_lower_first}StartTimeStr := c.DefaultQuery("${tag_name}_start", "")
    ${tag_name_lower_first}EndTimeStr := c.DefaultQuery("${tag_name}_end", "")

    var ${tag_name_lower_first}Start,${tag_name_lower_first}End time.Time
    if ${tag_name_lower_first}StartTimeStr != "" && ${tag_name_lower_first}EndTimeStr != "" {
        var e error
        ${tag_name_lower_first}Start, e = time.ParseInLocation("${list_layout}", ${tag_name_lower_first}StartTimeStr, ${time_zone})
        if e!=nil {
            c.JSON(400, gin.H{"message": e.Error()})
            return
        }
        ${tag_name_lower_first}End, e = time.ParseInLocation("${list_layout}", ${tag_name_lower_first}EndTimeStr, ${time_zone})
        if e!=nil {
            c.JSON(400, gin.H{"message": e.Error()})
            return
        }
        engine = engine.Where("${tag_name} between ? and ?", ${tag_name_lower_first}Start,  ${tag_name_lower_first}End.AddDate(0, 0, 1))
    }
`
				tmp = strings.Replace(tmp, "${tag_name_lower_first}", LowerFirstLetter(UnderLineToHump(tagName)), -1)
				tmp = strings.Replace(tmp, "${layout}", layout, -1)
				tmp = strings.Replace(tmp, "${list_layout}", listLayout, -1)

				tmp = strings.Replace(tmp, "${time_zone}", timeZone, -1)

				tmp = strings.Replace(tmp, "${field_name}", fieldName, -1)
				tmp = strings.Replace(tmp, "${tag_name}", tagName, -1)
				tmp = strings.Replace(tmp, "${zero_value}", zeroValue, -1)

				result += tmp

				continue
			}
			tagName := HumpToUnderLine(vType.Field(i).Name)
			fieldName := LowerFistLetter(vType.Field(i).Name)

			// jump ignored fields
			if in(tagName, strings.Split(replacement[0]["${jump_fields}"], ",")) {
				continue
			}
			tmp := `
    ${field_name} := c.DefaultQuery("${tag_name}", "")
    if ${field_name} != "" {
        engine = engine.Where("${tag_name} = ?", ${field_name})
    }
`
			tmp = strings.Replace(tmp, "${field_name}", fieldName, -1)
			tmp = strings.Replace(tmp, "${tag_name}", tagName, -1)
			result += tmp
		}
		return result
	}()
	var result string
	result = strings.Replace(format, "${handle}", handle, -1)
	return result
}

// using errno 和 errmsg
func GenerateListWhereV2(src interface{}, withListArgs bool, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	var format = `
   ${handle}
`
	vValue := reflect.ValueOf(src)
	vType := reflect.TypeOf(src)
	handle := func() string {
		var result string
		for i := 0; i < vValue.NumField(); i++ {
			// only handle basic types, otherwise continue
			if !in(vType.Field(i).Type.String(), []string{
				"int", "int8", "int16", "int32", "int64",
				"float32", "float64",
				"string",
				"uint8", "uint16", "uint32", "uint64",
				"time.Time",
			}) {
				continue
			}

			// handle time.Time alone
			if vType.Field(i).Type.AssignableTo(reflect.TypeOf(time.Time{})) {
				tagName := HumpToUnderLine(vType.Field(i).Name)
				fieldName := LowerFistLetter(vType.Field(i).Name)
				zeroValue := GetZeroValue(vValue.Field(i).Interface())
				layout := replacement[0]["layout"]
				if layout == "" {
					layout = "2006-01-02 15:04:05"
				}
				listLayout := replacement[0]["${list_layout}"]
				if layout == "" {
					layout = "2006-01-02"
				}
				timeZone := replacement[0]["time_zone"]
				if timeZone == "" {
					timeZone = "time.Local"
				}
				tmp := `
    ${tag_name_lower_first}StartTimeStr := c.DefaultQuery("${tag_name}_start", "")
    ${tag_name_lower_first}EndTimeStr := c.DefaultQuery("${tag_name}_end", "")

    var ${tag_name_lower_first}Start,${tag_name_lower_first}End time.Time
    if ${tag_name_lower_first}StartTimeStr != "" && ${tag_name_lower_first}EndTimeStr != "" {
        var e error
        ${tag_name_lower_first}Start, e = time.ParseInLocation("${list_layout}", ${tag_name_lower_first}StartTimeStr, ${time_zone})
        if e!=nil {
            c.JSON(400, gin.H{"errmsg": e.Error(), "errno":-1})
            return
        }
        ${tag_name_lower_first}End, e = time.ParseInLocation("${list_layout}", ${tag_name_lower_first}EndTimeStr, ${time_zone})
        if e!=nil {
            c.JSON(400, gin.H{"errmsg": e.Error(), "errno":-1})
            return
        }
        engine = engine.Where("${tag_name} between ? and ?", ${tag_name_lower_first}Start,  ${tag_name_lower_first}End.AddDate(0, 0, 1))
    }
`
				tmp = strings.Replace(tmp, "${tag_name_lower_first}", LowerFirstLetter(UnderLineToHump(tagName)), -1)
				tmp = strings.Replace(tmp, "${layout}", layout, -1)
				tmp = strings.Replace(tmp, "${list_layout}", listLayout, -1)

				tmp = strings.Replace(tmp, "${time_zone}", timeZone, -1)

				tmp = strings.Replace(tmp, "${field_name}", fieldName, -1)
				tmp = strings.Replace(tmp, "${tag_name}", tagName, -1)
				tmp = strings.Replace(tmp, "${zero_value}", zeroValue, -1)

				result += tmp

				continue
			}
			tagName := HumpToUnderLine(vType.Field(i).Name)
			fieldName := LowerFistLetter(vType.Field(i).Name)

			// jump ignored fields
			if in(tagName, strings.Split(replacement[0]["${jump_fields}"], ",")) {
				continue
			}
			tmp := `
    ${field_name} := c.DefaultQuery("${tag_name}", "")
    if ${field_name} != "" {
        engine = engine.Where("${tag_name} = ?", ${field_name})
    }
`
			tmp = strings.Replace(tmp, "${field_name}", fieldName, -1)
			tmp = strings.Replace(tmp, "${tag_name}", tagName, -1)
			result += tmp
		}
		return result
	}()
	var result string
	result = strings.Replace(format, "${handle}", handle, -1)
	return result
}

// Generate list api code.
// To completely use these code, you might import:
// "github.com/fwhezfwhez/errorx"
// mc_util "github.com/model_convert/util"
// you can get 'errorx.Wrap(e)','util.ToLimitOffset()', 'util.GenerateOrderBy()' above
//
// Replacement optional as:
// - ${page} "page"
// - ${size} "size"
// - ${order_by} ""
// - ${util_pkg} "util"
// - ${db_instance} "db.DB"
// - ${handler_name} "HTTPListUser"
// - ${model} "model.User"
// - ${handle_error} "fmt.Println(e)"
// - ${jump_fields}, "password,pw"
// - ${layout}, "2006-01-02 15:04:03"
// - ${time_zone} "time.Local"
func GenerateListAPI(src interface{}, withListArgs bool, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}

	if replacement[0]["${page}"] == "" {
		replacement[0]["${page}"] = "page"
	}

	if replacement[0]["${size}"] == "" {
		replacement[0]["${size}"] = "size"
	}

	if replacement[0]["${db_instance}"] == "" {
		replacement[0]["${db_instance}"] = "db.DB"
	}

	if replacement[0]["${util_pkg}"] == "" {
		replacement[0]["${util_pkg}"] = "mc_util"
	}
	if replacement[0]["${list_layout}"] == "" {
		replacement[0]["${list_layout}"] = "2006-01-02"
	}

	if replacement[0]["${handler_name}"] == "" {
		replacement[0]["${handler_name}"] = "HTTPListUser"
	}
	if replacement[0]["${list_layout}"] == "" {
		replacement[0]["${list_layout}"] = "2006-01-02"
	}

	vType := reflect.TypeOf(src)
	if replacement[0]["${model}"] == "" {
		replacement[0]["${model}"] = vType.String()
	}

	if replacement[0]["${handle_error}"] == "" {
		replacement[0]["${handle_error}"] = "log.Println(e)"
	}
	var copyMap = make(map[string]string)
	for k, v := range replacement[0] {
		copyMap[k] = v
	}

	queryArgsStatement := GenerateListWhere(src, withListArgs, []map[string]string{copyMap}...)
	commonStatementf := `
    page := c.DefaultQuery("${page}", "1")
    size := c.DefaultQuery("${size}", "20")
    orderBy := c.DefaultQuery("order_by", "${order_by}")
    
    var count int
    if e:= engine.Count(&count).Error; e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    var list = make([]${model}, 0, 20)
    if count == 0 {
        c.JSON(200, gin.H{"message": "success", "count": 0, "data": list})
        return
    }

    limit, offset := ${util_pkg}.ToLimitOffset(size, page, count)
    engine = engine.Limit(limit).Offset(offset)

    if orderBy != "" {
        engine = engine.Order(${util_pkg}.GenerateOrderBy(orderBy))
    }

    if e:= engine.Find(&list).Error; e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success", "count": count, "data": list})
`
	var resultf, result string
	resultf = queryArgsStatement + commonStatementf
	result = strings.Replace(resultf, "${model}", replacement[0]["${model}"], -1)

	result = strings.Replace(result, "${page}", replacement[0]["${page}"], -1)
	result = strings.Replace(result, "${size}", replacement[0]["${size}"], -1)
	result = strings.Replace(result, "${order_by}", replacement[0]["${order_by}"], -1)

	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	result = strings.Replace(result, "${util_pkg}", replacement[0]["${util_pkg}"], -1)

	var tmpf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateListAPI().
func ${handler_name}(c *gin.Context) {
    var engine = ${db_instance}.Model(&${model}{})
    ${result}
}
`
	tmp := strings.Replace(tmpf, "${handler_name}", replacement[0]["${handler_name}"], -1)
	log.Println(replacement[0]["db_instance"])
	log.Println(replacement[0])
	log.Println(replacement[0]["db_instance"] == "db.DB")
	tmp = strings.Replace(tmp, "${db_instance}", replacement[0]["${db_instance}"], -1)
	tmp = strings.Replace(tmp, "${model}", replacement[0]["${model}"], -1)
	tmp = strings.Replace(tmp, "${result}", result, -1)

	// format
	tmp = strings.Replace(tmp, "\n    \n", "\n", -1)
	tmp = strings.Replace(tmp, "\n\n", "\n", -1)
	tmp = strings.Replace(tmp, "\n\n}", "\n}", -1)
	tmp = strings.Replace(tmp, "\n    \n}", "\n}", -1)
	return tmp
}

// errmsg errno组合
func GenerateListAPIV2(src interface{}, withListArgs bool, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}

	if replacement[0]["${page}"] == "" {
		replacement[0]["${page}"] = "page"
	}

	if replacement[0]["${size}"] == "" {
		replacement[0]["${size}"] = "size"
	}

	if replacement[0]["${db_instance}"] == "" {
		replacement[0]["${db_instance}"] = "db.DB"
	}

	if replacement[0]["${util_pkg}"] == "" {
		replacement[0]["${util_pkg}"] = "mc_util"
	}
	if replacement[0]["${list_layout}"] == "" {
		replacement[0]["${list_layout}"] = "2006-01-02"
	}

	if replacement[0]["${handler_name}"] == "" {
		replacement[0]["${handler_name}"] = "HTTPListUser"
	}
	if replacement[0]["${list_layout}"] == "" {
		replacement[0]["${list_layout}"] = "2006-01-02"
	}

	vType := reflect.TypeOf(src)
	if replacement[0]["${model}"] == "" {
		replacement[0]["${model}"] = vType.String()
	}

	if replacement[0]["${handle_error}"] == "" {
		replacement[0]["${handle_error}"] = "log.Println(e)"
	}
	var copyMap = make(map[string]string)
	for k, v := range replacement[0] {
		copyMap[k] = v
	}

	queryArgsStatement := GenerateListWhereV2(src, withListArgs, []map[string]string{copyMap}...)
	commonStatementf := `
    page := c.DefaultQuery("${page}", "1")
    size := c.DefaultQuery("${size}", "20")
    orderBy := c.DefaultQuery("order_by", "${order_by}")
    
    var count int
    if e:= engine.Count(&count).Error; e!=nil {
        ${handle_error}
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }
    var list = make([]${model}, 0, 20)
    if count == 0 {
        c.JSON(200, gin.H{"errmsg": "success", "count": 0, "data": list, "errno":0})
        return
    }

    limit, offset := ${util_pkg}.ToLimitOffset(size, page, count)
    engine = engine.Limit(limit).Offset(offset)

    if orderBy != "" {
        engine = engine.Order(${util_pkg}.GenerateOrderBy(orderBy))
    }

    if e:= engine.Find(&list).Error; e!=nil {
        ${handle_error}
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }
    c.JSON(200, gin.H{"errmsg": "success", "count": count, "data": list, "errno":0})
`
	var resultf, result string
	resultf = queryArgsStatement + commonStatementf
	result = strings.Replace(resultf, "${model}", replacement[0]["${model}"], -1)

	result = strings.Replace(result, "${page}", replacement[0]["${page}"], -1)
	result = strings.Replace(result, "${size}", replacement[0]["${size}"], -1)
	result = strings.Replace(result, "${order_by}", replacement[0]["${order_by}"], -1)

	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	result = strings.Replace(result, "${util_pkg}", replacement[0]["${util_pkg}"], -1)

	var tmpf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateListAPI().
func ${handler_name}(c *gin.Context) {
    var engine = ${db_instance}.Model(&${model}{})
    ${result}
}
`
	tmp := strings.Replace(tmpf, "${handler_name}", replacement[0]["${handler_name}"], -1)
	log.Println(replacement[0]["db_instance"])
	log.Println(replacement[0])
	log.Println(replacement[0]["db_instance"] == "db.DB")
	tmp = strings.Replace(tmp, "${db_instance}", replacement[0]["${db_instance}"], -1)
	tmp = strings.Replace(tmp, "${model}", replacement[0]["${model}"], -1)
	tmp = strings.Replace(tmp, "${result}", result, -1)

	// format
	tmp = strings.Replace(tmp, "\n    \n", "\n", -1)
	tmp = strings.Replace(tmp, "\n\n", "\n", -1)
	tmp = strings.Replace(tmp, "\n\n}", "\n}", -1)
	tmp = strings.Replace(tmp, "\n    \n}", "\n}", -1)
	return tmp
}

// Generate get one api code.
// To completely use these code, you might import:
// "github.com/fwhezfwhez/errorx"
// you can get 'errorx.Wrap(e)' above
//
// - ${db_instance} "db.DB"
// - ${handler_name} "HTTPListUser"
// - ${model} "model.User"
// - ${handle_error} "fmt.Println(e, string(debug.Stack()))"
func GenerateGetOneAPI(src interface{}, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	handleDefault(src, replacement[0])
	var resultf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateGetOneAPI().
func ${handler_name}(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
        return
    }
    var count int
    if e:= (${db_instance}.Model(&${model}{}).Where("id=?", idInt).Count(&count).Error); e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
        return
    }
    var instance ${model}
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", id).First(&instance).Error); e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success", "data": instance})
}
`
	result := strings.Replace(resultf, "${handler_name}", replacement[0]["${handler_name}"], -1)
	result = strings.Replace(result, "${db_instance}", replacement[0]["${db_instance}"], -1)
	result = strings.Replace(result, "${model}", replacement[0]["${model}"], -1)
	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	return result
}

func GenerateGetOneAPIV2(src interface{}, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	handleDefault(src, replacement[0])
	var resultf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateGetOneAPI().
func ${handler_name}(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(200, gin.H{"errmsg": fmt.Sprintf("param 'id' requires int but got %s", id), "errno":-1})
        return
    }
    var count int
    if e:= (${db_instance}.Model(&${model}{}).Where("id=?", idInt).Count(&count).Error); e!=nil {
        ${handle_error}
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"errmsg": fmt.Sprintf("id '%s' record not found", id), "errno":-1})
        return
    }
    var instance ${model}
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", id).First(&instance).Error); e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }
    c.JSON(200, gin.H{"errmsg": "success", "data": instance, "errno":-1})
}
`
	result := strings.Replace(resultf, "${handler_name}", replacement[0]["${handler_name}"], -1)
	result = strings.Replace(result, "${db_instance}", replacement[0]["${db_instance}"], -1)
	result = strings.Replace(result, "${model}", replacement[0]["${model}"], -1)
	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	return result
}


// Generate add one api code.
// To completely use these code, you might import:
// "github.com/fwhezfwhez/errorx"
// you can get 'errorx.Wrap(e)' above
//
// - ${db_instance} "db.DB"
// - ${handler_name} "HTTPListUser"
// - ${model} "model.User"
// - ${handle_error} "fmt.Println(e, string(debug.Stack()))"
// - ${redis_conn} "redistool.RedisPool.Get()"
func GenerateAddOneAPI(src interface{}, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	handleDefault(src, replacement[0])

	var resultf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateAddOneAPI().
func ${handler_name} (c *gin.Context) {
    var param ${model}
    if e := c.Bind(&param); e!=nil {
        c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }

    if e:=(${db_instance}.Model(&${model}{}).Create(&param).Error); e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }

    if param.RedisKey() != "" {
        conn := ${redis_conn}
        defer conn.Close()
        param.DeleteFromRedis(conn)
    }
    c.JSON(200, gin.H{"message": "success", "data": param})
}
`
	result := strings.Replace(resultf, "${handler_name}", replacement[0]["${handler_name}"], -1)
	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	result = strings.Replace(result, "${db_instance}", replacement[0]["${db_instance}"], -1)
	result = strings.Replace(result, "${model}", replacement[0]["${model}"], -1)
	result = strings.Replace(result, "${redis_conn}", replacement[0]["${redis_conn}"], -1)
	return result
}

// 和V1相比，返回值组合为errno和errmsg
func GenerateAddOneAPIV2(src interface{}, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	handleDefault(src, replacement[0])

	var resultf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateAddOneAPI().
func ${handler_name} (c *gin.Context) {
    var param ${model}
    if e := c.Bind(&param); e!=nil {
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }

    if e:=(${db_instance}.Model(&${model}{}).Create(&param).Error); e!=nil {
        ${handle_error}
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }

    if param.RedisKey() != "" {
        conn := ${redis_conn}
        defer conn.Close()
        param.DeleteFromRedis(conn)
    }
    c.JSON(200, gin.H{"errmsg": "success", "data": param, "errno":0})
}
`
	result := strings.Replace(resultf, "${handler_name}", replacement[0]["${handler_name}"], -1)
	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	result = strings.Replace(result, "${db_instance}", replacement[0]["${db_instance}"], -1)
	result = strings.Replace(result, "${model}", replacement[0]["${model}"], -1)
	result = strings.Replace(result, "${redis_conn}", replacement[0]["${redis_conn}"], -1)
	return result
}

// Generate delete one api code.
// To completely use these code, you might import:
// "github.com/fwhezfwhez/errorx"
// you can get 'errorx.Wrap(e)' above
//
// - ${db_instance} "db.DB"
// - ${handler_name} "HTTPListUser"
// - ${model} "model.User"
// - ${handle_error} "fmt.Println(e, string(debug.Stack()))"
// - ${redis_conn} "redistool.RedisPool.Get()"
func GenerateDeleteOneAPI(src interface{}, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	handleDefault(src, replacement[0])
	var resultf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateDeleteOneAPI().
func ${handler_name}(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
        return
    }
    var count int
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", idInt).Count(&count).Error); e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
        return
    }
    var instance ${model}
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", idInt).First(&instance).Error); e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }    
    
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", id).Delete(&${model}{}).Error); e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if instance.RedisKey() != "" {
        conn := ${redis_conn}
        defer conn.Close()
        instance.DeleteFromRedis(conn)
    }
    c.JSON(200, gin.H{"message": "success"})
}
`
	result := strings.Replace(resultf, "${handler_name}", replacement[0]["${handler_name}"], -1)
	result = strings.Replace(result, "${db_instance}", replacement[0]["${db_instance}"], -1)
	result = strings.Replace(result, "${model}", replacement[0]["${model}"], -1)
	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	result = strings.Replace(result, "${redis_conn}", replacement[0]["${redis_conn}"], -1)
	return result
}

func GenerateDeleteOneAPIV2(src interface{}, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	handleDefault(src, replacement[0])
	var resultf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateDeleteOneAPI().
func ${handler_name}(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(200, gin.H{"errmsg": fmt.Sprintf("param 'id' requires int but got %s", id), "errno":-1})
        return
    }
    var count int
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", idInt).Count(&count).Error); e!=nil {
        ${handle_error}
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"errmsg": fmt.Sprintf("id '%s' record not found", id), "errno":-1})
        return
    }
    var instance ${model}
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", idInt).First(&instance).Error); e!=nil {
        ${handle_error}
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }    
    
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", id).Delete(&${model}{}).Error); e!=nil {
        ${handle_error}
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }
    if instance.RedisKey() != "" {
        conn := ${redis_conn}
        defer conn.Close()
        instance.DeleteFromRedis(conn)
    }
    c.JSON(200, gin.H{"errmsg": "success", "errno":0})
}
`
	result := strings.Replace(resultf, "${handler_name}", replacement[0]["${handler_name}"], -1)
	result = strings.Replace(result, "${db_instance}", replacement[0]["${db_instance}"], -1)
	result = strings.Replace(result, "${model}", replacement[0]["${model}"], -1)
	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	result = strings.Replace(result, "${redis_conn}", replacement[0]["${redis_conn}"], -1)
	return result
}

// Generate update one api code.
// To completely use these code, you might import:
// "github.com/fwhezfwhez/errorx"
// you can get 'errorx.Wrap(e)' above
//
// | field optional | default value | example value
// - ${util_pkg} | "util" | "model_convert.util"
// - ${args_forbid_update} | "" | "user_id, game_id"
// - ${db_instance} "db.DB" | "db.DB"
// - ${handler_name} "HTTPListUser" | HTTPUpdateUser |
// - ${model} "model.User" | "payModel.Order"
// - ${handle_error} "fmt.Println(e, string(debug.Stack()))" | raven.Throw(e)
// - ${redis_conn} "redistool.RedisPool.Get()"
func GenerateUpdateOneAPI(src interface{}, replacement ...map[string]string) string {
	return generateUpdateOneAPI(src, nil, replacement...)
}

func generateUpdateOneAPI(src interface{}, rangement []Field, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	handleDefault(src, replacement[0])
	var argsForbidf = `
    if !util.IfZero(param.${field_name}) {
        c.JSON(400, gin.H{"message": "field '${field_name}' can't be modified'"})
        return
    }
`
	var argsForbid string
	if replacement[0]["${args_forbid_update}"] != "" {
		forbids := strings.Split(replacement[0]["${args_forbid_update}"], ",")
		for _, arg := range forbids {
			arg = strings.TrimSpace(arg)
			arg = UnderLineToHump(arg)
			arg = UpperFirstLetter(arg)
			argsForbid += strings.Replace(argsForbidf, "${field_name}", arg, -1)
		}
	}

	var resultf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateUpdateOneAPI().
func ${handler_name}(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
        return
    }
    var count int
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", idInt).Count(&count).Error); e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
        return
    }

    var param ${model}
    if e:=c.Bind(&param);e!=nil {
        c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    
    ${args_forbid}
    
    var instance ${model}
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", idInt).First(&instance).Error); e!=nil {
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
	
    tran := ${db_instance}.Begin()

    // handle none-zero
    if e := tran.Model(&${model}{}).Where("id=?", id).Updates(param).Error; e!=nil {
        tran.Rollback()
        ${handle_error}
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
	
    // handle zero
	${zero_rangement}

    tran.Commit()

    if instance.RedisKey()!=""{
        conn := ${redis_conn}
        defer conn.Close()
        instance.DeleteFromRedis(conn)
    }
    c.JSON(200, gin.H{"message": "success"})
}
`
	var zeroRangement = `
    var m =  make(map[string]interface{})
		
    ${fill_m}

    if len(m) >0 {
        if e := tran.Model(&${model}{}).Where("id=?", id).Updates(m).Error; e!=nil {
            tran.Rollback()
            ${handle_error}
            c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
            return
        }
    }
	`

	var fillM string
	for _, v := range rangement {
		if In(v.TagName, []string{"id", "ID"}) {
			continue
		}
		if !In(v.TypeName, []string{
			"int", "int64", "int32", "int16", "int8", "uint8", "uint16", "uint32", "uint64", "byte", "float", "float64", "float32", "float16",
			"string",
		}) {
			continue
		}
		var tmpf string
		switch v.TypeName {
		case "string":
			tmpf = `
    if param.${field_name} == "12306" {
        m["${tag_name}"] = ""
    }

			`
		case "int", "int64", "int32", "int16", "int8", "uint8", "uint16", "uint32", "uint64", "byte", "float", "float64", "float32", "float16":
			tmpf = `
    if param.${field_name} == 12306 {
        m["${tag_name}"] = 0
    }

			`
		}

		tmp := strings.Replace(tmpf, "${field_name}", v.FieldName, -1)
		tmp = strings.Replace(tmp, "${tag_name}", v.TagName, -1)
		fillM += tmp
	}
	zeroRangement = strings.Replace(zeroRangement, "${fill_m}", fillM, -1)

	result := strings.Replace(resultf, "${handler_name}", replacement[0]["${handler_name}"], -1)

	result = strings.Replace(result, "${zero_rangement}", zeroRangement, -1)

	result = strings.Replace(result, "${db_instance}", replacement[0]["${db_instance}"], -1)
	result = strings.Replace(result, "${model}", replacement[0]["${model}"], -1)
	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	result = strings.Replace(result, "${redis_conn}", replacement[0]["${redis_conn}"], -1)

	result = strings.Replace(result, "${args_forbid}", argsForbid, -1)
	result = Format(result)
	return result
}

func generateUpdateOneAPIV2(src interface{}, rangement []Field, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	handleDefault(src, replacement[0])
	var argsForbidf = `
    if !util.IfZero(param.${field_name}) {
        c.JSON(200, gin.H{"errmsg": "field '${field_name}' can't be modified'", "errno":-1})
        return
    }
`
	var argsForbid string
	if replacement[0]["${args_forbid_update}"] != "" {
		forbids := strings.Split(replacement[0]["${args_forbid_update}"], ",")
		for _, arg := range forbids {
			arg = strings.TrimSpace(arg)
			arg = UnderLineToHump(arg)
			arg = UpperFirstLetter(arg)
			argsForbid += strings.Replace(argsForbidf, "${field_name}", arg, -1)
		}
	}

	var resultf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateUpdateOneAPI().
func ${handler_name}(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(200, gin.H{"errmsg": fmt.Sprintf("param 'id' requires int but got %s", id), "errno":-1})
        return
    }
    var count int
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", idInt).Count(&count).Error); e!=nil {
        ${handle_error}
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"errmsg": fmt.Sprintf("id '%s' record not found", id), "errno":-1})
        return
    }

    var param ${model}
    if e:=c.Bind(&param);e!=nil {
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }
    
    ${args_forbid}
    
    var instance ${model}
    if e:=(${db_instance}.Model(&${model}{}).Where("id=?", idInt).First(&instance).Error); e!=nil {
        ${handle_error}
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }
	
    tran := ${db_instance}.Begin()

    // handle none-zero
    if e := tran.Model(&${model}{}).Where("id=?", id).Updates(param).Error; e!=nil {
        tran.Rollback()
        ${handle_error}
        c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":-1})
        return
    }
	
    // handle zero
	${zero_rangement}

    tran.Commit()

    if instance.RedisKey()!=""{
        conn := ${redis_conn}
        defer conn.Close()
        instance.DeleteFromRedis(conn)
    }
    c.JSON(200, gin.H{"errmsg": "success", "errno":0})
}
`
	var zeroRangement = `
    var m =  make(map[string]interface{})
		
    ${fill_m}

    if len(m) >0 {
        if e := tran.Model(&${model}{}).Where("id=?", id).Updates(m).Error; e!=nil {
            tran.Rollback()
            ${handle_error}
            c.JSON(200, gin.H{"errmsg": errorx.Wrap(e).Error(), "errno":1})
            return
        }
    }
	`

	var fillM string
	for _, v := range rangement {
		if In(v.TagName, []string{"id", "ID"}) {
			continue
		}
		if !In(v.TypeName, []string{
			"int", "int64", "int32", "int16", "int8", "uint8", "uint16", "uint32", "uint64", "byte", "float", "float64", "float32", "float16",
			"string",
		}) {
			continue
		}
		var tmpf string
		switch v.TypeName {
		case "string":
			tmpf = `
    if param.${field_name} == "12306" {
        m["${tag_name}"] = ""
    }

			`
		case "int", "int64", "int32", "int16", "int8", "uint8", "uint16", "uint32", "uint64", "byte", "float", "float64", "float32", "float16":
			tmpf = `
    if param.${field_name} == 12306 {
        m["${tag_name}"] = 0
    }

			`
		}

		tmp := strings.Replace(tmpf, "${field_name}", v.FieldName, -1)
		tmp = strings.Replace(tmp, "${tag_name}", v.TagName, -1)
		fillM += tmp
	}
	zeroRangement = strings.Replace(zeroRangement, "${fill_m}", fillM, -1)

	result := strings.Replace(resultf, "${handler_name}", replacement[0]["${handler_name}"], -1)

	result = strings.Replace(result, "${zero_rangement}", zeroRangement, -1)

	result = strings.Replace(result, "${db_instance}", replacement[0]["${db_instance}"], -1)
	result = strings.Replace(result, "${model}", replacement[0]["${model}"], -1)
	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	result = strings.Replace(result, "${redis_conn}", replacement[0]["${redis_conn}"], -1)

	result = strings.Replace(result, "${args_forbid}", argsForbid, -1)
	result = Format(result)
	return result
}


func GenerateCacheProfileAPI(src interface{}, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}
	handleDefault(src, replacement[0])
	var resultf = `
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateCacheProfileAPI().
// You might import /path/to/${pkg_name_prefix}
func ${handler_name}(c *gin.Context) {
	type Param struct {
		Command string ${command_json_tag}
		Key     string ${key_json_tag}
	}
	var param Param
	if e := c.Bind(&param); e != nil {
		c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
		return
	}
	if param.Command == "" && param.Key == "" {
		c.JSON(200, gin.H{
			"message":            "success",
			"cache_name":         "${struct_name}",
			"cache_length":       len(${model}Cache),
			"cache_order_length": len(${model}CacheKeyOrder),
			"cache_config": gin.H{
				"delete_rate":  ${model}DeleteRate,
				"max_length":   ${pkg_name_prefix}${struct_name}CacheMaxLength,
				"cache_switch": ${pkg_name_prefix}${struct_name}CacheSwitch,
				"key_format": ${model}RedisKeyFormat,
			},

			"array_cache_name":         "Array${struct_name}",
			"array_cache_length":       len(${pkg_name_prefix}Array${struct_name}Cache),
			"array_cache_order_length": len(${pkg_name_prefix}Array${struct_name}CacheKeyOrder),
			"array_cache_config": gin.H{
				"delete_rate":  ${pkg_name_prefix}Array${struct_name}DeleteRate,
				"max_length":   ${pkg_name_prefix}Array${struct_name}CacheMaxLength,
				"cache_switch": ${pkg_name_prefix}Array${struct_name}CacheSwitch,
				"array_key_format": ${pkg_name_prefix}Array${struct_name}RedisKeyFormat,
			},
		})
		return
	}

	if param.Key != "" && param.Command == "" {
		${model}CacheLock.RLock()
		defer ${model}CacheLock.RUnlock()
		v, ok := ${model}Cache[param.Key]
		if ok {
			c.JSON(200, gin.H{
				"message": "success",
				"tip":     "hit in ${model}Cache",
				"key":     param.Key,
				"value":   v,
			})
			return
		}
		${pkg_name_prefix}Array${struct_name}CacheLock.RLock()
		defer ${pkg_name_prefix}Array${struct_name}CacheLock.RUnlock()
		v2, ok2 := ${pkg_name_prefix}Array${struct_name}Cache[param.Key]
		if ok2 {
			c.JSON(200, gin.H{
				"message": "success",
				"tip":     "hit in Array${struct_name}Cache",
				"key":     param.Key,
				"value":   v2,
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "success",
			"tip":     "hit none nether in ${model}Cache nor Array${struct_name}Cache",
			"key":     param.Key,
			"value":   v2,
		})
		return
	}
	c.JSON(200, gin.H{"message": "success", "tip": "I'm confused, maybe you can put nothing in body"})
}
`
	result := strings.Replace(resultf, "${handler_name}", replacement[0]["${handler_name}"], -1)
	result = strings.Replace(result, "${command_json_tag}", "`json:\"commond\"`", -1)
	result = strings.Replace(result, "${key_json_tag}", "`json:\"key\"`", -1)
	result = strings.Replace(result, "${pkg_name_prefix}", replacement[0]["${pkg_name_prefix}"], -1)

	result = strings.Replace(result, "${db_instance}", replacement[0]["${db_instance}"], -1)
	result = strings.Replace(result, "${model}", replacement[0]["${model}"], -1)
	result = strings.Replace(result, "${struct_name}", replacement[0]["${struct_name}"], -1)

	result = strings.Replace(result, "${handle_error}", replacement[0]["${handle_error}"], -1)
	result = strings.Replace(result, "${redis_conn}", replacement[0]["${redis_conn}"], -1)
	result = Format(result)
	return result
}

type Field struct {
	isZero    bool
	TypeName  string
	FieldName string
	TagName   string
	Value     interface{}
}

func GenerateCRUD(src interface{}, replacement ...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}

	var rangement = make([]Field, 0, 10)
	vType := reflect.TypeOf(src)
	vValue := reflect.ValueOf(src)
	for i := 0; i < vType.NumField(); i++ {
		tagStr := vType.Field(i).Tag.Get("json")
		if tagStr == "-" || tagStr == "" {
			continue
		}
		arr := strings.Split(tagStr, ",")
		tagValue := strings.TrimSpace(arr[0])
		valueI := vValue.Field(i).Interface()
		rangement = append(rangement, Field{
			isZero:    IfZero(valueI),
			TypeName:  vType.Field(i).Type.Name(),
			FieldName: vType.Field(i).Name,
			TagName:   tagValue,
			Value:     valueI,
		})
	}

	modelName := vType.Name()

	var rs string
	replacement[0]["${handler_name}"] = "HTTPAdd" + modelName
	addAPI := GenerateAddOneAPI(src, replacement...)

	replacement[0]["${handler_name}"] = "HTTPList" + modelName
	listAPI := GenerateListAPI(src, false, replacement...)

	replacement[0]["${handler_name}"] = "HTTPGet" + modelName
	getAPI := GenerateGetOneAPI(src, replacement...)

	replacement[0]["${handler_name}"] = "HTTPUpdate" + modelName
	updateAPI := generateUpdateOneAPI(src, rangement, replacement...)

	replacement[0]["${handler_name}"] = "HTTPDelete" + modelName
	deleteAPI := GenerateDeleteOneAPI(src, replacement...)

	//replacement[0]["${handler_name}"] = "HTTPProfileCache" + modelName
	//profileCacheAPI := GenerateCacheProfileAPI(src, replacement...)

	var note = `
// Auto generated by github.com/fwhezfwhez/model_convert.GenerateCRUD. You might need import:
// "github.com/gin-gonic/gin"
// "github.com/fwhezfwhez/errorx"
// "github.com/fwhezfwhez/model_convert/util"
//
// "package/path/to/${db_instance}"
// "package/path/to/${model}"
// "package/path/to/${redis_conn}"
//
// Use codes below like:
/*
	r := gin.Default()
	r.POST("/${url_letters}/", ${generate_to_pkg}.HTTPAdd${struct_name})
	r.PATCH("/${url_letters}/:id/", ${generate_to_pkg}.HTTPUpdate${struct_name})
	r.DELETE("/${url_letters}/:id/", ${generate_to_pkg}.HTTPDelete${struct_name})
	r.GET("/${url_letters}/", ${generate_to_pkg}.HTTPList${struct_name})
	r.GET("/${url_letters}/:id/", ${generate_to_pkg}.HTTPGet${struct_name})
*/
`
	note = strings.Replace(note, "${db_instance}", replacement[0]["${db_instance}"], -1)
	note = strings.Replace(note, "${model}", replacement[0]["${model}"], -1)
	note = strings.Replace(note, "${struct_name}", replacement[0]["${struct_name}"], -1)
	note = strings.Replace(note, "${generate_to_pkg}", replacement[0]["${generate_to_pkg}"], -1)
	note = strings.Replace(note, "${url_letters}", URLLetter(replacement[0]["${struct_name}"]), -1)

	rs = fmt.Sprintf("%s%s  %s  %s  %s  %s", note, addAPI, listAPI, getAPI, updateAPI, deleteAPI)
	rs = Format(rs)

	// 注入剪贴板
	clipboard.WriteAll(rs)
	return rs
}

func handleDefault(src interface{}, replacement map[string]string) {

	if replacement["${page}"] == "" {
		replacement["${page}"] = "page"
	}

	if replacement["${size}"] == "" {
		replacement["${size}"] = "size"
	}

	if replacement["${db_instance}"] == "" {
		replacement["${db_instance}"] = "db.DB"
	}

	if replacement["${util_pkg}"] == "" {
		replacement["${util_pkg}"] = "util"
	}

	if replacement["${handler_name}"] == "" {
		replacement["${handler_name}"] = "HTTPListUser"
	}

	vType := reflect.TypeOf(src)
	if replacement["${model}"] == "" {
		replacement["${model}"] = vType.String()
	}
	arr := strings.Split(replacement["${model}"], ".")
	if len(arr) == 2 {
		replacement["${pkg_name_prefix}"] = arr[0] + "."
		replacement["${struct_name}"] = arr[1]
	} else if len(arr) == 1 {
		replacement["${struct_name}"] = arr[0]
	} else {
		// do nothing
	}

	if replacement["${generate_to_pkg}"] == "" {
		replacement["${generate_to_pkg}"] = strings.ToLower(replacement["struct_name"])
	}
	if replacement["${handle_error}"] == "" {
		replacement["${handle_error}"] = "fmt.Println(e, string(debug.Stack()))"
	}
	if replacement["${redis_conn}"] == "" {
		replacement["${redis_conn}"] = "redistool.RedisPool.Get()"
	}
}

// 和V1相比，这个版本，支持errmsg，errno
func GenerateCRUDV2(src interface{}, replacement...map[string]string) string {
	if len(replacement) == 0 {
		replacement = []map[string]string{
			map[string]string{},
		}
	}

	var rangement = make([]Field, 0, 10)
	vType := reflect.TypeOf(src)
	vValue := reflect.ValueOf(src)
	for i := 0; i < vType.NumField(); i++ {
		tagStr := vType.Field(i).Tag.Get("json")
		if tagStr == "-" || tagStr == "" {
			continue
		}
		arr := strings.Split(tagStr, ",")
		tagValue := strings.TrimSpace(arr[0])
		valueI := vValue.Field(i).Interface()
		rangement = append(rangement, Field{
			isZero:    IfZero(valueI),
			TypeName:  vType.Field(i).Type.Name(),
			FieldName: vType.Field(i).Name,
			TagName:   tagValue,
			Value:     valueI,
		})
	}

	modelName := vType.Name()

	var rs string
	replacement[0]["${handler_name}"] = "HTTPAdd" + modelName
	addAPI := GenerateAddOneAPIV2(src, replacement...)

	replacement[0]["${handler_name}"] = "HTTPList" + modelName
	listAPI := GenerateListAPIV2(src, false, replacement...)

	replacement[0]["${handler_name}"] = "HTTPGet" + modelName
	getAPI := GenerateGetOneAPIV2(src, replacement...)

	replacement[0]["${handler_name}"] = "HTTPUpdate" + modelName
	updateAPI := generateUpdateOneAPIV2(src, rangement, replacement...)

	replacement[0]["${handler_name}"] = "HTTPDelete" + modelName
	deleteAPI := GenerateDeleteOneAPIV2(src, replacement...)

	//replacement[0]["${handler_name}"] = "HTTPProfileCache" + modelName
	//profileCacheAPI := GenerateCacheProfileAPI(src, replacement...)

	var note = `
// Auto generated by github.com/fwhezfwhez/model_convert.GenerateCRUD. You might need import:
// "github.com/gin-gonic/gin"
// "github.com/fwhezfwhez/errorx"
// "github.com/fwhezfwhez/model_convert/util"
//
// "package/path/to/${db_instance}"
// "package/path/to/${model}"
// "package/path/to/${redis_conn}"
//
// Use codes below like:
/*
	r := gin.Default()
	r.POST("/${url_letters}/", ${generate_to_pkg}.HTTPAdd${struct_name})
	r.PATCH("/${url_letters}/:id/", ${generate_to_pkg}.HTTPUpdate${struct_name})
	r.DELETE("/${url_letters}/:id/", ${generate_to_pkg}.HTTPDelete${struct_name})
	r.GET("/${url_letters}/", ${generate_to_pkg}.HTTPList${struct_name})
	r.GET("/${url_letters}/:id/", ${generate_to_pkg}.HTTPGet${struct_name})
*/
`
	note = strings.Replace(note, "${db_instance}", replacement[0]["${db_instance}"], -1)
	note = strings.Replace(note, "${model}", replacement[0]["${model}"], -1)
	note = strings.Replace(note, "${struct_name}", replacement[0]["${struct_name}"], -1)
	note = strings.Replace(note, "${generate_to_pkg}", replacement[0]["${generate_to_pkg}"], -1)
	note = strings.Replace(note, "${url_letters}", URLLetter(replacement[0]["${struct_name}"]), -1)

	rs = fmt.Sprintf("%s%s  %s  %s  %s  %s", note, addAPI, listAPI, getAPI, updateAPI, deleteAPI)
	rs = Format(rs)

	// 注入剪贴板
	clipboard.WriteAll(rs)
	return rs
}
