package model_convert

func GetModelFromTable(dataSource string, tableName string) string {
	return TableToStructWithTag(dataSource, tableName)
}
