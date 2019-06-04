package model_convert

import (
	"fmt"
	"strings"
)

func XMLToModel(src string, structName string) string {
	src = strings.Trim(src, "\n")

	var result = `
       type %s struct{
           %s
       }
    `
	var xmlName string
	var lineTmp string
	lines := strings.Split(src, "\n")
	for _, v := range lines {
		v = strings.Trim(v, " ")
		v = strings.Trim(v, "\n")
		j := strings.Index(v, ">")
		if j == -1 {
			continue
		}
		if NumberTimesOf(v, "<") == 1 && xmlName == "" {
			// xml头
			j = strings.Index(v, ">")
			xmlName = string(v[1:j])
			lineTmp += fmt.Sprintf("XMLName xml.Name `xml:\"%s\"`\n", xmlName)
			continue
		} else {
			// 表体
			if string(v[1]) == "/" {
				// 表尾
				continue
			}
			if strings.Contains(v, "CDATA") {
				lineTmp += fmt.Sprintf("            %s string `xml:\"%s,CDATA\"`\n", UnderLineToHump(v[1:j]), string(v[1:j]))
			} else {
				lineTmp += fmt.Sprintf("            %s string `xml:\"%s\"`\n", UnderLineToHump(v[1:j]), string(v[1:j]))
			}
		}
	}
	result = fmt.Sprintf(result, structName, lineTmp)
	return result
}

func NumberTimesOf(str string, single string) int {
	var times int
	for i, _ := range str {
		if string(str[i]) == single {
			times++
		}
	}
	return times
}
