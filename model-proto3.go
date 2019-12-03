package model_convert

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// i is a go model instance ,not ptr.
// replacement should only length by 1 and it will specific some field as setting. There are optional replace-able field:
// map[string]string {
//     "${pb_pkg_name}": "userProto", // default pb
//     "${model_pkg_name}": "userModel",  // default model
// }
//
// If optional fields value is 'empty' then it's empty.If do not set 'empty' will use its defaults:
// "${pb_pkg_name}": "empty"
//
func GoModelToProto3(i interface{}, replacement ... map[string]string) (string, string, string) {
	if len(replacement) > 1 {
		panic("model_convert.GoModelToProto3 only requries replacement length less by 1 but got " + strconv.Itoa(len(replacement)))
	}

	var pbPkgName string
	var modelPkgName string
	if replacement != nil {
		pbPkgName = replacement[0]["${pb_pkg_name}"]
		modelPkgName = replacement[0]["${model_pkg_name}"]
	}
	if pbPkgName == "" {
		pbPkgName = "pb"
	}
	if modelPkgName == "" {
		modelPkgName = "model"
	}

	if replacement[0]["${pb_pkg_name}"] == "empty" {
		pbPkgName = ""
	}
	if replacement[0]["${model_pkg_name}"] == "empty" {
		modelPkgName = ""
	}

	var t []uint8
	// Get message definition
	var resultF = "message %s {\n%s}"
	vType := reflect.TypeOf(i)
	vValue := reflect.ValueOf(i)

	var middle string
	var pIndex int
	for i := 0; i < vValue.NumField(); i++ {
		tmp := vType.Field(i).Type.AssignableTo(reflect.TypeOf(t))
		if vType.Field(i).Type.Kind() == reflect.Slice && tmp != true {
			var linef = "    repeated %s %s =%d;\n"
			var line string
			var pType string
			var pName string

			tmp2 := vType.Field(i).Type.String()
			_ = tmp2
			pType = typeConvertFromGoToProto3(vType.Field(i).Type.String())
			pName = HumpToUnderLine(vType.Field(i).Name)
			pIndex++
			line = fmt.Sprintf(linef, pType, pName, pIndex)
			middle += line

			continue
		}

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
	var SetMf string

	SetMf = `
func SetModel%s(src ${pb_pkg_name}.%s) ${model_pkg_name}.%s {
    var dest ${model_pkg_name}.%s
    %s
    return dest
}
    `

	var SetM string

	if pbPkgName == "" {
		SetM = strings.Replace(SetMf, "${pb_pkg_name}.", pbPkgName, -1)
	} else {
		SetM = strings.Replace(SetMf, "${pb_pkg_name}", pbPkgName, -1)
	}

	if modelPkgName == "" {
		SetM = strings.Replace(SetM, "${model_pkg_name}.", modelPkgName, -1)
	} else {
		SetM = strings.Replace(SetM, "${model_pkg_name}", modelPkgName, -1)
	}
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
	var SetPf = `
func SetProto%s(src ${model_pkg_name}.%s) ${pb_pkg_name}.%s {
    var dest ${pb_pkg_name}.%s
    %s
    return dest
}
    `

	var SetP string

	if pbPkgName == "" {
		SetP = strings.Replace(SetPf, "${pb_pkg_name}.", pbPkgName, -1)
	} else {
		SetP = strings.Replace(SetPf, "${pb_pkg_name}", pbPkgName, -1)
	}

	if modelPkgName == "" {
		SetP = strings.Replace(SetP, "${model_pkg_name}.", modelPkgName, -1)
	} else {
		SetP = strings.Replace(SetP, "${model_pkg_name}", modelPkgName, -1)
	}

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

// i is a go model instance ,not ptr.
// replacement should only length by 1 and it will specific some field as setting. There are optional replace-able field:
// map[string]string {
//     "${pb_pkg_name}": "userProto", // default pb
//     "${model_pkg_name}": "userModel",  // default model
//     "${start_index}": "1", // default 0
// }
//
// If optional fields value is 'empty' then it's empty.If do not set 'empty' will use its defaults:
// "${pb_pkg_name}": "empty"
//
func GoModelToProto2(i interface{}, replacement ... map[string]string) (string, string, string) {
	if len(replacement) > 1 {
		panic("model_convert.GoModelToProto3 only requries replacement length less by 1 but got " + strconv.Itoa(len(replacement)))
	}

	var pbPkgName string
	var modelPkgName string
	if replacement != nil {
		pbPkgName = replacement[0]["${pb_pkg_name}"]
		modelPkgName = replacement[0]["${model_pkg_name}"]
	}
	if pbPkgName == "" {
		pbPkgName = "pb"
	}
	if modelPkgName == "" {
		modelPkgName = "model"
	}

	if replacement[0]["${pb_pkg_name}"] == "empty" {
		pbPkgName = ""
	}
	if replacement[0]["${model_pkg_name}"] == "empty" {
		modelPkgName = ""
	}

	var t []uint8
	// Get message definition
	var resultF = "message %s {\n%s}"
	vType := reflect.TypeOf(i)
	vValue := reflect.ValueOf(i)

	var middle string
	var pIndex int

	if replacement[0]["${start_index}"] != "" {
		pIndex, _ = strconv.Atoi(replacement[0]["${start_index}"])
	}

	// go model, key fieldName, value golang type name
	var typeMap = make(map[string]string)

	// go model, key type name, value proto transfer
	var protoTypeMap = map[string]string{
		"int":        "proto.Int32",
		"string":     "proto.String",
		"RawMessage": "[]byte",
	}

	for i := 0; i < vValue.NumField(); i++ {
		tmp := vType.Field(i).Type.AssignableTo(reflect.TypeOf(t))
		if vType.Field(i).Type.Kind() == reflect.Slice && tmp != true {
			var linef = "    optional repeated %s %s =%d;\n"
			var line string
			var pType string
			var pName string

			tmp2 := vType.Field(i).Type.String()
			_ = tmp2
			pType = typeConvertFromGoToProto3(vType.Field(i).Type.String())
			pName = HumpToUnderLine(vType.Field(i).Name)
			pIndex++
			line = fmt.Sprintf(linef, pType, pName, pIndex)
			middle += line

			continue
		}

		typeMap[vType.Field(i).Name] = vType.Field(i).Type.Name()

		var linef = "    optional %s %s =%d;\n"
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
	var SetMf string

	SetMf = `
func SetModel%s(src ${pb_pkg_name}.%s) ${model_pkg_name}.%s {
    var dest ${model_pkg_name}.%s
    %s
    return dest
}
    `

	var SetM string

	if pbPkgName == "" {
		SetM = strings.Replace(SetMf, "${pb_pkg_name}.", pbPkgName, -1)
	} else {
		SetM = strings.Replace(SetMf, "${pb_pkg_name}", pbPkgName, -1)
	}

	if modelPkgName == "" {
		SetM = strings.Replace(SetM, "${model_pkg_name}.", modelPkgName, -1)
	} else {
		SetM = strings.Replace(SetM, "${model_pkg_name}", modelPkgName, -1)
	}
	var line2 string
	var suffix = "\n    "
	for i := 0; i < vValue.NumField(); i++ {
		if i == vValue.NumField()-1 {
			suffix = ""
		}
		var linef = "dest.%s = %s(src.Get%s())" + suffix
		var propName = vType.Field(i).Name
		line2 += fmt.Sprintf(linef, propName, typeMap[propName], propName)
	}

	// Set go struct function
	var SetPf = `
func SetProto%s(src ${model_pkg_name}.%s) ${pb_pkg_name}.%s {
    var dest ${pb_pkg_name}.%s
    %s
    return dest
}
    `

	var SetP string

	if pbPkgName == "" {
		SetP = strings.Replace(SetPf, "${pb_pkg_name}.", pbPkgName, -1)
	} else {
		SetP = strings.Replace(SetPf, "${pb_pkg_name}", pbPkgName, -1)
	}

	if modelPkgName == "" {
		SetP = strings.Replace(SetP, "${model_pkg_name}.", modelPkgName, -1)
	} else {
		SetP = strings.Replace(SetP, "${model_pkg_name}", modelPkgName, -1)
	}

	var line3 string
	suffix = "\n    "

	for i := 0; i < vValue.NumField(); i++ {
		if i == vValue.NumField()-1 {
			suffix = ""
		}
		var linef = "dest.%s = %s(int32(src.%s))" + suffix
		var propName = vType.Field(i).Name
		line3 += fmt.Sprintf(linef, propName, protoTypeMap[typeMap[propName]], propName)
	}

	var setM = fmt.Sprintf(SetM, vType.Name(), vType.Name(), vType.Name(), vType.Name(), line2)
	var setP = fmt.Sprintf(SetP, vType.Name(), vType.Name(), vType.Name(), vType.Name(), line3)
	var message = fmt.Sprintf(resultF, vType.Name(), middle)
	return message, setM, setP
}

func typeConvertFromGoToProto3(gtype string) string {
	var ptype string
	switch gtype {
	case "string":
		ptype = "string"
	case "json.RawMessage", "[]byte", "[]uint8":
		ptype = "bytes"
	case "int", "int32", "int8", "uint8", "uint16", "uint32":
		ptype = "int32"
	case "int64":
		ptype = "int64"
	default:
		tmp := strings.Split(gtype, ".")
		ptype = tmp[len(tmp)-1]
	}
	return ptype
}
