package model_convert

import (
	"fmt"
	"reflect"
)

func GoModelToProto3(i interface{}) (string, string, string) {

	// Get message definition
	var resultF = "message %s {\n%s}"
	vType := reflect.TypeOf(i)
	vValue := reflect.ValueOf(i)

	var middle string
	var pIndex int
	for i := 0; i < vValue.NumField(); i++ {
		var linef = "    %s %s =%d;\n"
		var line string
		var pType string
		var pName string

		pType = typeConvertFromGoToProto3(vType.Field(i).Type.String())
		pName = HumpToUnderLine(vType.Field(i).Name)
		pIndex++
		line = fmt.Sprintf(linef, pType, pName, pIndex)
		middle += line
	}

	// Set proto struct function
	var SetM = `
func SetModel%s(src pb.%s) model.%s {
    var dest model.%s
    %s
    return dest
}
    `
	var line2 string
	var suffix = "\n    "
	for i := 0; i < vValue.NumField(); i++ {
		if i == vValue.NumField()-1 {
			suffix = ""
		}
		var linef = "dest.%s = src.%s" + suffix
		var propName = vType.Field(i).Name
		line2 += fmt.Sprintf(linef, propName, propName)
	}

	// Set go struct function
	var SetP = `
func SetProto%s(src model.%s) pb.%s {
    var dest pb.%s
    %s
    return dest
}
    `
	var line3 string
	suffix = "\n    "
	for i := 0; i < vValue.NumField(); i++ {
		if i == vValue.NumField()-1 {
			suffix = ""
		}
		var linef = "dest.%s = src.%s" + suffix
		var propName = vType.Field(i).Name
		line3 += fmt.Sprintf(linef, propName, propName)
	}

	return fmt.Sprintf(resultF, vType.Name(), middle),
		fmt.Sprintf(SetM, vType.Name(), vType.Name(), vType.Name(), vType.Name(), line2),
		fmt.Sprintf(SetP, vType.Name(), vType.Name(), vType.Name(), vType.Name(), line3)

}

func typeConvertFromGoToProto3(gtype string) string {
	var ptype string
	switch gtype {
	case "string":
		ptype = "string"
	case "json.RawMessage", "[]byte":
		ptype = "bytes"
	case "int", "int32", "int8", "uint8", "uint16", "uint32":
		ptype = "int32"
	case "int64":
		ptype = "int64"
	}
	return ptype
}
