package model_convert

import (
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
    if startTimeStr != "" && endTimeStr !="" {
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
    if ${field_name}!= "" {
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
