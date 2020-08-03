package model_convert

import (
	"strings"
)

func GenerateMarkDown(src interface{}, context ...map[string]interface{}) string {
	if len(context) == 0 {
		context = []map[string]interface{}{
			map[string]interface{}{},
		}
	}
	var rs string
	add := generateMDAdd(src, context[0])
	update := generateMDUpdate(src, context[0])
	list := generateMDList(src, context[0])
	get := generateMDGet(src, context[0])
	del := generateMDDelete(src, context[0])

	rs = add + update + list + get + del

	return rs
}

func generateMDList(src interface{}, context ...map[string]interface{}) string {

	if len(context) == 0 {
		context = []map[string]interface{}{
			map[string]interface{}{},
		}
	}
	_, ok := context[0]["${sub_order}"]
	if !ok {
		context[0]["${sub_order}"] = 1
	} else {
		context[0]["${sub_order}"] = context[0]["${sub_order}"].(int) + 1
	}

	context[0]["${progres}"] = "GET-LIST"
	defer func() {
		context[0]["${progres}"] = ""
	}()

	AnalysisSrc(src, context[0])
	// ${H1} ${md_order}. ${model_chinese_name}增删改查(CRUD)，后台配置
	rs := `
${H2} ${md_order}.${sub_order} 列表${model_chinese_name}

${backquote_3}url
GET /${model-name}/
${backquote_3}

* 所有字段均参与排序

${backquote_3}url
按照id降序排列: ?order_by=-id, 升序排列: ?order_by=id
按照created_at降序，id升序: ?order_by=-created_at,id
${backquote_3}

* 所有字段均参与筛选，对时间字段筛选时，为 ${backquote}?字段名_start=2020-01-01&字段名_end=2020-03-01${backquote}
* 默认每页20条，可以通过page指定页数，size指定每页条数 ${backquote}?page=1&size=20${backquote}

返回:

* status 200

${backquote_3}json
{
    "message": "success",
    "count": 1,
    "data":[${model_json_4}]
}
${backquote_3}

${model_table}

* status 400/500

${backquote_3}json
{
    "message":"出错信息"
}
${backquote_3}
`

	rs = strings.Replace(rs, "${backquote}", BackQuote(1), -1)
	rs = strings.Replace(rs, "${backquote_3}", BackQuote(3), -1)

	rs = strings.Replace(rs, "${H1}", MarkDownH(1), -1)
	rs = strings.Replace(rs, "${H2}", MarkDownH(2), -1)
	rs = strings.Replace(rs, "${H3}", MarkDownH(3), -1)
	rs = strings.Replace(rs, "${H4}", MarkDownH(4), -1)

	rs = ReplaceDefault(rs, "${sub_order}", context[0]["${sub_order}"], -1, "1")

	rs = ReplaceDefault(rs, "${model-name}", context[0]["${model-name}"], -1, "user-info")

	rs = ReplaceDefault(rs, "${md_order}", context[0]["${md_order}"], -1, "1")
	rs = ReplaceDefault(rs, "${model_chinese_name}", context[0]["${model_chinese_name}"], -1, "用户信息")
	rs = ReplaceDefault(rs, "${model_json}", context[0]["${model_json}"], -1, "{}")
	rs = ReplaceDefault(rs, "${model_json_no_id}", context[0]["${model_json_no_id}"], -1, "{}")

	rs = ReplaceDefault(rs, "${model_json_4}", context[0]["${model_json_4}"], -1, "{}")
	rs = ReplaceDefault(rs, "${model_table}", context[0]["${model_table}"], -1, "")

	return rs
}

func generateMDGet(src interface{}, context ...map[string]interface{}) string {

	if len(context) == 0 {
		context = []map[string]interface{}{
			map[string]interface{}{},
		}
	}
	_, ok := context[0]["${sub_order}"]
	if !ok {
		context[0]["${sub_order}"] = 1
	} else {
		context[0]["${sub_order}"] = context[0]["${sub_order}"].(int) + 1
	}

	context[0]["${progres}"] = "GET-LIST"
	defer func() {
		context[0]["${progres}"] = ""
	}()

	AnalysisSrc(src, context[0])
	// ${H1} ${md_order}. ${model_chinese_name}增删改查(CRUD)，后台配置
	rs := `
${H2} ${md_order}.${sub_order} 获取某个${model_chinese_name}

${backquote_3}url
GET /${model-name}/:id/
${backquote_3}


返回:

* status 200

${backquote_3}json
{
    "message": "success",
    "data":${model_json_4}
}
${backquote_3}

${model_table}

* status 400/500

${backquote_3}json
{
    "message":"出错信息"
}
${backquote_3}
`

	rs = strings.Replace(rs, "${backquote}", BackQuote(1), -1)
	rs = strings.Replace(rs, "${backquote_3}", BackQuote(3), -1)

	rs = strings.Replace(rs, "${H1}", MarkDownH(1), -1)
	rs = strings.Replace(rs, "${H2}", MarkDownH(2), -1)
	rs = strings.Replace(rs, "${H3}", MarkDownH(3), -1)
	rs = strings.Replace(rs, "${H4}", MarkDownH(4), -1)

	rs = ReplaceDefault(rs, "${sub_order}", context[0]["${sub_order}"], -1, "1")

	rs = ReplaceDefault(rs, "${model-name}", context[0]["${model-name}"], -1, "user-info")

	rs = ReplaceDefault(rs, "${md_order}", context[0]["${md_order}"], -1, "1")
	rs = ReplaceDefault(rs, "${model_chinese_name}", context[0]["${model_chinese_name}"], -1, "用户信息")
	rs = ReplaceDefault(rs, "${model_json}", context[0]["${model_json}"], -1, "{}")
	rs = ReplaceDefault(rs, "${model_json_no_id}", context[0]["${model_json_no_id}"], -1, "{}")

	rs = ReplaceDefault(rs, "${model_json_4}", context[0]["${model_json_4}"], -1, "{}")
	rs = ReplaceDefault(rs, "${model_table}", context[0]["${model_table}"], -1, "")

	return rs
}

func generateMDAdd(src interface{}, context ...map[string]interface{}) string {

	if len(context) == 0 {
		context = []map[string]interface{}{
			map[string]interface{}{},
		}
	}

	_, ok := context[0]["${sub_order}"]
	if !ok {
		context[0]["${sub_order}"] = 1
	} else {
		context[0]["${sub_order}"] = context[0]["${sub_order}"].(int) + 1
	}

	context[0]["${progres}"] = "POST-ADD"
	defer func() {
		context[0]["${progres}"] = ""
	}()

	AnalysisSrc(src, context[0])

	rs := `
${H1} ${md_order}. ${model_chinese_name}增删改查(CRUD)，后台配置
- 本文档由脚本自动生成，仅限制字段的命名和类型。具体限制，需要以实际为主。

${H2} ${md_order}.${sub_order} 新增${model_chinese_name}

${backquote_3}url
POST /${model-name}/
Content-Type application/json
${backquote_3}

body:

${backquote_3}json
${model_json_no_id}
${backquote_3}

${model_table}

返回:

* status 200

${backquote_3}json
{
    "message": "success",
    "data": ${model_json_4}
}
${backquote_3}

* status 400/500

${backquote_3}json
{
    "message":"出错信息"
}
${backquote_3}
`

	rs = strings.Replace(rs, "${backquote}", BackQuote(1), -1)
	rs = strings.Replace(rs, "${backquote_3}", BackQuote(3), -1)

	rs = strings.Replace(rs, "${H1}", MarkDownH(1), -1)
	rs = strings.Replace(rs, "${H2}", MarkDownH(2), -1)
	rs = strings.Replace(rs, "${H3}", MarkDownH(3), -1)
	rs = strings.Replace(rs, "${H4}", MarkDownH(4), -1)

	rs = ReplaceDefault(rs, "${sub_order}", context[0]["${sub_order}"], -1, "1")

	rs = ReplaceDefault(rs, "${model-name}", context[0]["${model-name}"], -1, "user-info")

	rs = ReplaceDefault(rs, "${md_order}", context[0]["${md_order}"], -1, "1")
	rs = ReplaceDefault(rs, "${model_chinese_name}", context[0]["${model_chinese_name}"], -1, "用户信息")
	rs = ReplaceDefault(rs, "${model_json}", context[0]["${model_json}"], -1, "{}")
	rs = ReplaceDefault(rs, "${model_json_no_id}", context[0]["${model_json_no_id}"], -1, "{}")

	rs = ReplaceDefault(rs, "${model_json_4}", context[0]["${model_json_4}"], -1, "{}")
	rs = ReplaceDefault(rs, "${model_table}", context[0]["${model_table}"], -1, "")

	return rs
}

func generateMDUpdate(src interface{}, context ...map[string]interface{}) string {
	if len(context) == 0 {
		context = []map[string]interface{}{
			map[string]interface{}{},
		}
	}
	_, ok := context[0]["${sub_order}"]
	if !ok {
		context[0]["${sub_order}"] = 1
	} else {
		context[0]["${sub_order}"] = context[0]["${sub_order}"].(int) + 1
	}

	context[0]["${progres}"] = "PATCH-UPDATE"
	defer func() {
		context[0]["${progres}"] = ""
	}()

	AnalysisSrc(src, context[0])
	// ${H1} ${md_order}. ${model_chinese_name}增删改查(CRUD)，后台配置
	rs := `
${H2} ${md_order}.${sub_order} 修改${model_chinese_name}

${backquote_3}url
PATCH /${model-name}/:id/
Content-Type application/json
${backquote_3}

body:

${backquote_3}json
${model_json_no_id}
${backquote_3}

* 数字类型对0更新不敏感，字符串类型对 ${backquote}""${backquote} 更新不敏感，布尔类型对false不敏感
* 数字类型传12306时，表示修改对应字段为0，字符串类型传 ${backquote}""${backquote} 将修改为空字符串
* 只传某一个字段时，表示修改该字段，而不需要每个字段的原始值都传

${model_table}

返回:

* status 200

${backquote_3}json
{
    "message": "success"
}
${backquote_3}

* status 400/500

${backquote_3}json
{
    "message":"出错信息"
}
${backquote_3}
`

	rs = strings.Replace(rs, "${backquote}", BackQuote(1), -1)
	rs = strings.Replace(rs, "${backquote_3}", BackQuote(3), -1)

	rs = strings.Replace(rs, "${H1}", MarkDownH(1), -1)
	rs = strings.Replace(rs, "${H2}", MarkDownH(2), -1)
	rs = strings.Replace(rs, "${H3}", MarkDownH(3), -1)
	rs = strings.Replace(rs, "${H4}", MarkDownH(4), -1)

	rs = ReplaceDefault(rs, "${sub_order}", context[0]["${sub_order}"], -1, "1")

	rs = ReplaceDefault(rs, "${model-name}", context[0]["${model-name}"], -1, "user-info")

	rs = ReplaceDefault(rs, "${md_order}", context[0]["${md_order}"], -1, "1")
	rs = ReplaceDefault(rs, "${model_chinese_name}", context[0]["${model_chinese_name}"], -1, "用户信息")
	rs = ReplaceDefault(rs, "${model_json}", context[0]["${model_json}"], -1, "{}")
	rs = ReplaceDefault(rs, "${model_json_no_id}", context[0]["${model_json_no_id}"], -1, "{}")

	rs = ReplaceDefault(rs, "${model_json_4}", context[0]["${model_json_4}"], -1, "{}")
	rs = ReplaceDefault(rs, "${model_table}", context[0]["${model_table}"], -1, "")

	return rs
}

func generateMDDelete(src interface{}, context ...map[string]interface{}) string {
	if len(context) == 0 {
		context = []map[string]interface{}{
			map[string]interface{}{},
		}
	}
	_, ok := context[0]["${sub_order}"]
	if !ok {
		context[0]["${sub_order}"] = 1
	} else {
		context[0]["${sub_order}"] = context[0]["${sub_order}"].(int) + 1
	}

	context[0]["${progres}"] = "DELETE-DELETE"
	defer func() {
		context[0]["${progres}"] = ""
	}()

	AnalysisSrc(src, context[0])
	// ${H1} ${md_order}. ${model_chinese_name}增删改查(CRUD)，后台配置
	rs := `
${H2} ${md_order}.${sub_order} 删除${model_chinese_name}

${backquote_3}url
DELETE /${model-name}/:id/
Content-Type application/json
${backquote_3}


返回:

* status 200

${backquote_3}json
{
    "message": "success"
}
${backquote_3}

* status 400/500

${backquote_3}json
{
    "message":"出错信息"
}
${backquote_3}
`

	rs = strings.Replace(rs, "${backquote}", BackQuote(1), -1)
	rs = strings.Replace(rs, "${backquote_3}", BackQuote(3), -1)

	rs = strings.Replace(rs, "${H1}", MarkDownH(1), -1)
	rs = strings.Replace(rs, "${H2}", MarkDownH(2), -1)
	rs = strings.Replace(rs, "${H3}", MarkDownH(3), -1)
	rs = strings.Replace(rs, "${H4}", MarkDownH(4), -1)

	rs = ReplaceDefault(rs, "${sub_order}", context[0]["${sub_order}"], -1, "1")

	rs = ReplaceDefault(rs, "${model-name}", context[0]["${model-name}"], -1, "user-info")

	rs = ReplaceDefault(rs, "${md_order}", context[0]["${md_order}"], -1, "1")
	rs = ReplaceDefault(rs, "${model_chinese_name}", context[0]["${model_chinese_name}"], -1, "用户信息")
	rs = ReplaceDefault(rs, "${model_json}", context[0]["${model_json}"], -1, "{}")
	rs = ReplaceDefault(rs, "${model_json_no_id}", context[0]["${model_json_no_id}"], -1, "{}")

	rs = ReplaceDefault(rs, "${model_json_4}", context[0]["${model_json_4}"], -1, "{}")
	rs = ReplaceDefault(rs, "${model_table}", context[0]["${model_table}"], -1, "")

	return rs
}
