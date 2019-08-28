package model_convert

import (
	"reflect"
)

func Model2PGSQL(src interface{}) (string, string) {
	var sql = `
create table ${tableName}(
   ${fields}
)
`
	var reModel = `
type ${modelName} struct{
   ${fields}
}

func (o ${structName}) TableName() string{
    return ${tableName}
}

${handleUpperTypeFields}
   `
	vValue := reflect.ValueOf(src)
	vType := reflect.TypeOf(src)

	for i := 0; i < vValue.NumField(); i++ {
		_ = vType
	}
	return sql, reModel
}
