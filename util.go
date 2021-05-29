package model_convert

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fwhezfwhez/errorx"

	"github.com/shopspring/decimal"
)

func LowerFistLetter(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(string(s[0])) + s[1:]
}
func UpperFirstLetter(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
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
func URLLetter(s string) string {
	s1 := HumpToUnderLine(s)
	return strings.Join(strings.Split(s1, "_"), "-")
}

func GetZeroValue(src interface{}) string {
	switch src.(type) {
	case int, int8, int16, int32, int64, float32, float64, uint, uint8, uint16, uint32, uint64:
		return "0"
	case string:
		return `""`
	}
	return `""`
}

func Format(arg string) string {
	arg = strings.Replace(arg, "\n\n", "\n", -1)
	arg = strings.Replace(arg, "\n\n\n", "\n", -1)
	arg = strings.Replace(arg, "\n    \n", "\n", -1)
	arg = strings.Replace(arg, "//\n\n//", "//\n//", -1)
	return arg
}

func StringMaxLen(max int, buf []byte) string {
	if len(buf) == 0 {
		return ""
	}
	if max >= len(buf) {
		return string(buf)
	}
	return string(buf[:max])
}

func IfZero(arg interface{}) bool {
	if arg == nil {
		return true
	}
	switch v := arg.(type) {
	case int, int32, int16, int64:
		if v == 0 {
			return true
		}
	case float32:
		r := float64(v)
		return math.Abs(r-0) < 0.0000001
	case float64:
		return math.Abs(v-0) < 0.0000001
	case string:
		if v == "" || v == "%%" || v == "%" {
			return true
		}
	case *string, *int, *int64, *int32, *int16, *int8, *float32, *float64, *time.Time:
		if v == nil {
			return true
		}
	case time.Time:
		return v.IsZero()
	case decimal.Decimal:
		tmp, _ := v.Float64()
		return math.Abs(tmp-0) < 0.0000001
	default:
		return false
	}
	return false
}

func In(s string, arr []string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

// Output ` * num
func BackQuote(num int) string {
	var rs string
	for i := 0; i < num; i++ {
		rs += "`"
	}
	return rs
}
func MarkDownH(num int) string {
	var rs string
	for i := 0; i < num; i++ {
		rs += "#"
	}
	return rs
}

func AnalysisSrc(src interface{}, context map[string]interface{}) {
	// init context
	context["backquote-1"] = BackQuote(1)
	context["backquote-3"] = BackQuote(3)

	// ${rangement}
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

	context["${rangement}"] = rangement

	// ${model_name}, ${modelName}, ${model-name}
	arr := strings.Split(vType.Name(), ".")
	if len(arr) == 0 {
		fmt.Println(errorx.NewFromString("model_name not detected").Error())
		return
	}
	modelName := arr[len(arr)-1]
	context["${modelName}"] = LowerFirstLetter(modelName)
	context["${model_name}"] = HumpToUnderLine(modelName)
	context["${ModelName}"] = UpperFirstLetter(modelName)
	context["${model-name}"] = strings.Join(strings.Split(HumpToUnderLine(modelName), "_"), "-")

	// ${model_json}, ${model_json_4}, ${model_json_no_id}
	var m = make(map[string]interface{})
	for _, v := range rangement {
		m[HumpToUnderLine(v.TagName)] = randomValueOf(v.TypeName)
	}

	buf, e := json.MarshalIndent(m, "", "    ")
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		return
	}
	context["${model_json}"] = string(buf)
	buf, e = json.MarshalIndent(m, "    ", "    ")
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		return
	}
	context["${model_json_4}"] = string(buf)
	// context[0]["${progres}"] = "GET-LIST"
	_ = context["${progres}"]

	delete(m, "id")
	buf, e = json.MarshalIndent(m, "", "    ")
	if e != nil {
		fmt.Println(errorx.Wrap(e).Error())
		return
	}
	context["${model_json_no_id}"] = string(buf)

	// ${model_table}
	{
		var table = "\n|字段 | 类型 | 含义 | 示例 | 是否必要 |\n| :--- | :--- | :--- | :--- | :--- |\n"
		var progres = context["${progres}"]
		var isNeccessary string
		if progres == "POST-ADD" {
			isNeccessary = "是"
		} else if progres == "PATCH-UPDATE" {
			isNeccessary = "否"
		}
		for _, v := range rangement {

			// do not handle id to add an item
			if progres == "POST-ADD" || progres == "PATCH-UPDATE" {
				if v.TagName == "id" {
					continue
				}
			}

			line := fmt.Sprintf("| %s | %s | | | %s |\n", v.TagName, JSONType(v.TypeName), isNeccessary)
			table += line
		}
		context["${model_table}"] = table
	}
}

var randString = []string{"抗击肺炎", "China", "Golang", "hill", "hot-dog", "crazy", "中国江西", "江西中至", "zonst", "critical", "kill me heal me", "dot"}
var randNumber = []float64{1, 3.14, 15.1, 99, 397, 15, 29, 92812, 320931, 789, 2123.12, 82193.15, 91, 123, 12}
var randInt = []float64{1, 3, 15, 99, 397, 15, 29, 92812, 320931, 789, 2123, 82193, 91, 123, 12}

var randTime = []string{"2019-01-02 15:55:01", "2018-03-20 12:11:00", "2020-04:20 00:00"}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomValueOf(vType string) interface{} {
	num := r.Intn(10000)
	switch vType {
	case "string":
		return randString[num%len(randString)]
	case "time.Time", "Time":
		return randTime[num%len(randTime)]
	case "float64", "float32", "float":
		return randNumber[num%len(randTime)]
	case "int", "int8", "int16", "int32", "int64", "uint8", "uint", "uint16", "uint32", "uint64":
		return randInt[num%len(randInt)]
	}
	return fmt.Sprintf("undefine-vType-%s", vType)
}

func ReplaceDefault(s string, old string, new interface{}, n int, defaultValue string) string {
	var newStr string
	var newInt int
	var ok1, ok2 bool
	newInt, ok1 = new.(int)

	newStr, ok2 = new.(string)
	if !ok1 && !ok2 {
		newStr = defaultValue
	} else {
		if ok1 {
			newStr = strconv.Itoa(newInt)
		} else {
			// do nothing
		}
	}

	return strings.Replace(s, old, newStr, n)
}

func JSONType(vType string) string {
	switch vType {
	case "string":
		return "string"
	case "time.Time", "Time":
		return "string, datetime"
	case "int", "int8", "int16", "int32", "int64", "uint8", "uint", "uint16", "uint32", "uint64":
		return "int"
	case "float64", "float32", "float":
		return "number,float"
	}
	return fmt.Sprintf("undefine-json-type-mapping-%s", vType)
}

func Pick(rate float64) bool {
	// fmt.Println(r.Float64)

	// num := math.Ceil(1 / rate)
	// n := r.Intn(int(num))
	// return n < 1
	return r.Float64() <= rate
}

func GetDefault(i string, dft string) string {
	if i == "" {
		return dft
	}
	return i
}

func replaceAll(tmpl string, r map[string]interface{}) string {
	var rs string
	for k, v := range r {
		rs = strings.Replace(tmpl, k, fmt.Sprintf("%s", v), -1)
		tmpl = rs
	}
	return rs
}

func TrimLine(line string) string {
	line = strings.TrimSpace(line)
	line = strings.Trim(line, "\t")
	return line
}
