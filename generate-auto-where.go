package model_convert

import (
	"log"
	"reflect"
	"strings"
	"time"
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
func GenerateListWhere(src interface{}, withListArgs bool, replacement ... map[string]string) string {
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
	handle := func() (string) {
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
				timeZone := replacement[0]["time_zone"]
				if timeZone == "" {
					timeZone = "time.Local"
				}
				tmp := `
    startTimeStr := c.DefaultQuery("start_time", "")
    endTimeStr := c.DefaultQuery("end_time", "")

    var start,end time.Time
    if startTimeStr != "" && endTimeStr != "" {
        var e error
        start, e = time.ParseLocation("${layout}", startTimeStr, ${time_zone})
        if e!=nil {
            c.JSON(400, gin.H{"message": e.Error()})
            return
        }
        end, e = time.ParseLocation("${layout}", endTimeStr, ${time_zone})
        if e!=nil {
            c.JSON(400, gin.H{"message": e.Error()})
            return
        }
        engine = engine.Where("${tag_name} between ? and ?", start, end.Add(0, 0, 1))
    }
`
				tmp = strings.Replace(tmp, "${layout}", layout, -1)
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
        engine = engine.Where("${tag_name} != ?", ${field_name})
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
// "github.com/model_convert/util"
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
func GenerateListAPI(src interface{}, withListArgs bool, replacement ... map[string]string) string {
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
		replacement[0]["${util_pkg}"] = "util"
	}

	if replacement[0]["${handler_name}"] == "" {
		replacement[0]["${handler_name}"] = "HTTPListUser"
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
    c.JSON(200, gin.H{"message": "success", "count": 0, "data": list})
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
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateList().
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
