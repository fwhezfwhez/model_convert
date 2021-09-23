package model_convert

import (
	"github.com/atotto/clipboard"
	"strings"
)

func GenerateNote(sql string) string {
	var rsf = `
${input}
;
${notes}
`

	var input = sql

	notes := generateSqlNotes(sql)

	rs := replaceAll(rsf, map[string]interface{}{
		"${notes}": notes,
		"${input}": input,
	})

	// 注入剪贴板
	clipboard.WriteAll(rs)
	return rs
}

func generateSqlNotes(sql string) string {
	type Field struct {
		FieldName string
		Note      []string
	}

	lines := strings.Split(sql, "\n")
	if len(lines) == 0 {
		return sql
	}

	var fields = make([]Field, 0, 10)
	_ = fields

	var sumnotes []string
	var fieldFlag string
	_ = fieldFlag

	var tableName string

	for i, _ := range lines {
		line := TrimLine(lines[i])
		if len(line) == 0 {
			continue
		}

		if strings.Contains(line, "unique") {
			continue
		}

		if strings.Contains(line, "create table") {

			tmp := strings.TrimPrefix(line, "create table")
			tmp = strings.TrimSuffix(tmp, "(")

			tableName = TrimLine(tmp)

			// 获取表名
			continue
		}

		if len(line) == 1 && line == ")" {
			continue
		}

		if len(line) == 1 && line == ");" {
			continue
		}

		// 防止上一行注释被清理
		if len(sumnotes) == 0 {
			sumnotes = make([]string, 0, 10)
		}

		if strings.HasPrefix(line, "--") {
			sumnotes = append(sumnotes, strings.TrimPrefix(line, "--"))
		} else {
			fieldFlag = strings.Split(line, " ")[0]
		}

		// 抓取末尾的注释
		if strings.Contains(line, "--") && strings.Index(line, "--") != 0 {
			anote := strings.TrimPrefix(line[strings.Index(line, "--"):], "--")
			sumnotes = append(sumnotes, anote)
		}

		if fieldFlag != "" {

			field := Field{
				FieldName: fieldFlag,
				Note:      sumnotes,
			}

			if len(sumnotes) == 0 && field.FieldName == "updated_at" {
				field.Note = append(sumnotes, "更新于")
			}
			if len(sumnotes) == 0 && field.FieldName == "game_id" {
				field.Note = append(sumnotes, "平台id")
			}

			if len(sumnotes) == 0 && field.FieldName == "user_id" {
				field.Note = append(sumnotes, "用户id")
			}

			if len(sumnotes) == 0 && field.FieldName == "created_at" {
				field.Note = append(sumnotes, "创建于")
			}

			if len(sumnotes) == 0 && field.FieldName == "id" {
				field.Note = append(sumnotes, "id,主键")
			}

			fields = append(fields, field)

			sumnotes = nil
			fieldFlag = ""
		}
	}

	var format = "comment on column ${tableName}.${columnName} is '${note}';"

	var notes = make([]string, 0, 10)
	for _, v := range fields {
		note := replaceAll(format, map[string]interface{}{
			"${tableName}":  tableName,
			"${columnName}": v.FieldName,
			"${note}":       strings.Join(v.Note, ","),
		})

		notes = append(notes, note)
	}

	return strings.Join(notes, "\n")
}
