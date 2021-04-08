## model-convert
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/fwhezfwhez/model_convert)

model-convert is used for transfer all kinds of structs

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Declaration](#declaration)
- [1. Start](#1-start)
- [2. Stable cases](#2-stable-cases)
    - [2.1 Gorm-postgres/mysql table to model【supporting 1st-cache 2nd-cache】](#21-gorm-postgresmysql-table-to-modelsupporting-1st-cache-2nd-cache)
    - [2.2 Xml to go model](#22-xml-to-go-model)
    - [2.3 Add json,form tag for go model](#23-add-jsonform-tag-for-go-model)
    - [2.4 Go model transfer to protobuf3](#24-go-model-transfer-to-protobuf3)
    - [2.5 Http restful api](#25-http-restful-api)
        - [2.5.1 list](#251-list)
        - [2.5.2 get-one](#252-get-one)
        - [2.5.3 add-one](#253-add-one)
        - [2.5.4 delete-one](#254-delete-one)
        - [2.5.5 update-one](#255-update-one)
        - [2.5.6 generate crud](#256-generate-crud)
    - [2.6 grpc](#26-grpc)
      - [2.6.1 grpc-proto-service to go model](#261-grpc-proto-service-to-go-model)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## Declaration
** Most cases are supposed to use before app starts, they're not advised using in runtime.It helps auto-build structures.
** Most structure auto-built might need to be formatted yourself.

## 1. Start
`go get github.com/fwhezfwhez/model_convert`

## 2. Stable cases
Most cases are developing.However, only when specific requirements are met with, I will upgrade requiring functions. Here are stable use cases, it will be taken care of when project are updating.

#### 2.1 Gorm-postgres/mysql table to model【supporting 1st-cache 2nd-cache】

| usage | description | well practicing cases |
| --- | -- | -- |
| 1st-cache| Read from redis first, then access to db | All cases |
|2nd-cache | Read from cmap | config-only|

```go
package main

import (
    "fmt"
    "github.com/fwhezfwhez/model_convert"
)

func main() {
    // postgres
    dataSouce := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "localhost", "5432", "postgres", "game", "disable", "123")
    tableName := "user_info"
    fmt.Println(model_convert.TableToStructWithTag(dataSouce, tableName, "postgres"))

    // mysql
    // dataSouce = "ft:123@/test?charset=utf8&parseTime=True&loc=Local&allowNativePasswords=true"
    // tableName = "t_user"
    // fmt.Println(mc.TableToStructWithTag(dataSouce, tableName, "mysql"))
}
```
output:

postgres
```go
type UserInfo struct {
    Id        int    `gorm:"column:id;default:" json:"id" form:"id"`
    UserId    int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
    OpenId    string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`
    UnionId   string `gorm:"column:union_id;default:" json:"union_id" form:"union_id"`
    UserName  string `gorm:"column:user_name;default:" json:"user_name" form:"user_name"`
    HeaderUrl string `gorm:"column:header_url;default:" json:"header_url" form:"header_url"`
    Sex       int    `gorm:"column:sex;default:" json:"sex" form:"sex"`
    GameId    int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
}

func (o UserInfo) TableName() string {
    return "user_info"
}

// ... strenthening list:
// - 1.1st/2nd cache for single object/ objects array
// - 2.cmap 2nd-cache
```
mysql
```go
type TUser struct {
	Attach    json.RawMessage `gorm:"column:attach;default:" json:"attach" form:"attach"`
	CreatedAt time.Time       `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
	Id        int             `gorm:"column:id;default:" json:"id" form:"id"`
	Name      string          `gorm:"column:name;default:" json:"name" form:"name"`
}

func (o TUser) TableName() string {
	return "t_user"
}
```

#### 2.2 Xml to go model
```go
package main

import (
    "fmt"
    "github.com/fwhezfwhez/model_convert"
)

func main(){
    fmt.Println(model_convert.XMLToModel(`
        <xml>
            <ToUserName><![CDATA[toUser]]></ToUserName>
            <FromUserName><![CDATA[fromUser]]></FromUserName>
            <CreateTime>1348831860</CreateTime>
            <MsgType><![CDATA[text]]></MsgType>
            <Content><![CDATA[this is a test]]></Content>
            <MsgId>1234567890123456</MsgId>
        </xml>
    `, "MessageInfo"))
}
```
output:
```go
type MessageInfo struct{
    XMLName xml.Name `xml:"xml"`
    ToUserName   string `xml:"ToUserName,CDATA"`
    FromUserName string `xml:"FromUserName,CDATA"`
    CreateTime   string `xml:"CreateTime"`
    MsgType      string `xml:"MsgType,CDATA"`
    Content      string `xml:"Content,CDATA"`
    MsgId        string `xml:"MsgId"`
}
```

#### 2.3 Add json,form tag for go model
Only support under-line style `AaBb -> aa_bb`
```go
package main

import (
    "fmt"
    "github.com/fwhezfwhez/model_convert"
)


func main() {
    fmt.Println(model_convert.AddJSONFormTag(
        `
          type UserInfo struct {
            Id        int
            UserId    int
            OpenId    string
            UnionId   string
            UserName  string
            HeaderUrl string
            Sex       int
            GameId    int
        }
        `,
    ))
}

```
output:

```go
type UserInfo struct {
    Id        int    `json:"id" form:"id"`
    UserId    int    `json:"user_id" form:"user_id"`
    OpenId    string `json:"open_id" form:"open_id"`
    UnionId   string `json:"union_id" form:"union_id"`
    UserName  string `json:"user_name" form:"user_name"`
    HeaderUrl string `json:"header_url" form:"header_url"`
    Sex       int    `json:"sex" form:"sex"`
    GameId    int    `json:"game_id" form:"game_id"`
}
```
#### 2.4 Go model transfer to protobuf3
Developing.Requiring modify a bit by yourself.
```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/model_convert"
)

func main() {
	type U struct {
		Username string
		Password string
		Age      int
		Id       int32
		Config   json.RawMessage
	}
	ps, setM, setP := model_convert.GoModelToProto3(U{})
	fmt.Println(ps)
	fmt.Println(setM)
	fmt.Println(setP)
}

```

output:

```go
message U {
    string username =1;
    string password =2;
    int32 age =3;
    int32 id =4;
    bytes config =5;
}

func SetModelU(src pb.U) model.U {
    var dest model.U
    dest.Username = src.Username
    dest.Password = src.Password
    dest.Age = src.Age
    dest.Id = src.Id
    dest.Config = src.Config
    return dest
}


func SetProtoU(src model.U) pb.U {
    var dest pb.U
    dest.Username = src.Username
    dest.Password = src.Password
    dest.Age = src.Age
    dest.Id = src.Id
    dest.Config = src.Config
    return dest
}
```
#### 2.5 Http restful api
This provides auto generate api crud http code.Supporting only gin+gorm.

###### 2.5.1 list
Note:
```go
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
func model_convert.GenerateListAPI()
```

```go
package main
import (
  	"encoding/json"
  	"fmt"
  	"github.com/fwhezfwhez/model_convert"
)

func main() {
	type LingqianOrder struct {
		Id       int             `gorm:"column:id;default:" json:"id" form:"id"`
		OrderId  string          `gorm:"column:order_id;default:" json:"order_id" form:"order_id"`
		GameId   int             `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId   int             `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId   string          `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`
		State    int             `gorm:"column:state;default:" json:"state" form:"state"`
		Request  json.RawMessage `gorm:"column:request;default:" json:"request" form:"request"`
		Response string          `gorm:"column:response;default:" json:"response" form:"response"`
	}

	fmt.Println(model_convert.GenerateListAPI(LingqianOrder{}, false, map[string]string{
		"${model}": "payModel.LingqianOrder",
		"${handler_name}" : "HTTPListLingqianOrder",
		"${handle_error}": `common.SaveError(e)`,
	}))
}
```
Output
```go
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateList().
func HTTPListLingqianOrder(c *gin.Context) {
    var engine = db.DB.Model(&payModel.LingqianOrder{})

    id := c.DefaultQuery("id", "")
    if id != "" {
        engine = engine.Where("id != ?", id)
    }
    orderId := c.DefaultQuery("order_id", "")
    if orderId != "" {
        engine = engine.Where("order_id = ?", orderId)
    }
    gameId := c.DefaultQuery("game_id", "")
    if gameId != "" {
        engine = engine.Where("game_id = ?", gameId)
    }
    userId := c.DefaultQuery("user_id", "")
    if userId != "" {
        engine = engine.Where("user_id = ?", userId)
    }
    openId := c.DefaultQuery("open_id", "")
    if openId != "" {
        engine = engine.Where("open_id = ?", openId)
    }
    state := c.DefaultQuery("state", "")
    if state != "" {
        engine = engine.Where("state = ?", state)
    }
    response := c.DefaultQuery("response", "")
    if response != "" {
        engine = engine.Where("response = ?", response)
    }

    page := c.DefaultQuery("page", "1")
    size := c.DefaultQuery("size", "20")
    orderBy := c.DefaultQuery("order_by", "")
    var count int
    if e:= engine.Count(&count).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    var list = make([]payModel.LingqianOrder, 0, 20)
    if count == 0 {
        c.JSON(200, gin.H{"message": "success", "count": 0, "data": list})
        return
    }
    limit, offset := util.ToLimitOffset(size, page, count)
    engine = engine.Limit(limit).Offset(offset)
    if orderBy != "" {
        engine = engine.Order(util.GenerateOrderBy(orderBy))
    }
    if e:= engine.Find(&list).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success", "count": 0, "data": list})
}
```

###### 2.5.2 get-one
Note:
```go
// Generate get-one api code.
// To completely use these code, you might import:
// "github.com/fwhezfwhez/errorx"
// you can get 'errorx.Wrap(e)' above
//
// - ${db_instance} "db.DB"
// - ${handler_name} "HTTPListUser"
// - ${model} "model.User"
// - ${handle_error} "fmt.Println(e)"
func GenerateGetOneAPI()
```
Usage:
```go
package main
import (
    "encoding/json"
    "fmt"
    "github.com/fwhezfwhez/model_convert"
)
func main() {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int    `gorm:"column:state;default:" json:"state" form:"state"`
	}

	rs := model_convert.GenerateGetOneAPI(VxTemplateUser{}, map[string]string{
		"${model}": "payModel.LingqianOrder",
		"${handler_name}" : "HTTPGetOneLingqianOrder",
		"${handle_error}": `common.SaveError(e)`,
	})
	fmt.Println(rs)
}
```
Output:
```go
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateGetOneAPI().
func HTTPGetOneLingqianOrder(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
        return
    }
    var count int
    if e:=db.DB.Model(&payModel.LingqianOrder{}).Where("id=?", idInt).Count(&count).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
        return
    }
    var instance payModel.LingqianOrder
    if e:=db.DB.Model(&payModel.LingqianOrder{}).Where("id=?", id).First(&instance).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success", "data": instance})
}

```

###### 2.5.3 add-one
Note:
```go
// Generate add one api code.
// To completely use these code, you might import:
// "github.com/fwhezfwhez/errorx"
// you can get 'errorx.Wrap(e)' above
//
// - ${db_instance} "db.DB"
// - ${handler_name} "HTTPListUser"
// - ${model} "model.User"
// - ${handle_error} "fmt.Println(e, string(debug.Stack()))"
func GenerateAddOneAPI()
```
```go
package main
import (
    "encoding/json"
    "fmt"
    "github.com/fwhezfwhez/model_convert"
)
func main() {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int    `gorm:"column:state;default:" json:"state" form:"state"`
	}

	rs := model_convert.GenerateAddOneAPI(VxTemplateUser{}, map[string]string{
		"${model}": "payModel.LingqianOrder",
		"${handler_name}" : "HTTPAddLingqianOrder",
		"${handle_error}": `common.SaveError(e)`,
	})
	fmt.Println(rs)
}
```
Output:
```go
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateAddOneAPI().
func HTTPAddLingqianOrder (c *gin.Context) {
    var param payModel.LingqianOrder
    if e := c.Bind(&param); e!=nil {
        c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }

    if e:=db.DB.Model(&payModel.LingqianOrder{}).Create(&param).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success", "data": param})
}
```

###### 2.5.4 delete-one

Note:
```go
// Generate delete one api code.
// To completely use these code, you might import:
// "github.com/fwhezfwhez/errorx"
// you can get 'errorx.Wrap(e)' above
//
// - ${db_instance} "db.DB"
// - ${handler_name} "HTTPListUser"
// - ${model} "model.User"
// - ${handle_error} "fmt.Println(e, string(debug.Stack()))"
func GenerateDeleteOneAPI()
```
```go
package main
import (
    "encoding/json"
    "fmt"
    "github.com/fwhezfwhez/model_convert"
)
func main() {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int    `gorm:"column:state;default:" json:"state" form:"state"`
	}

	rs := GenerateDeleteOneAPI(VxTemplateUser{}, map[string]string{
		"${model}": "payModel.LingqianOrder",
		"${handler_name}" : "HTTPDeleteLingqianOrder",
		"${handle_error}": `common.SaveError(e)`,
	})
	fmt.Println(rs)
}
```
Output:
```go
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateDeleteOneAPI().
func HTTPDeleteLingqianOrder(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
        return
    }
    var count int
    if e:=db.DB.Model(&payModel.LingqianOrder{}).Where("id=?", idInt).Count(&count).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
        return
    }
    var instance payModel.LingqianOrder
    if e:=db.DB.Model(&payModel.LingqianOrder{}).Where("id=?", id).Delete(&instance).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success"})
}
```

###### 2.5.5 update-one

Note:
```go
// Generate update one api code.
// To completely use these code, you might import:
// "github.com/fwhezfwhez/errorx"
// you can get 'errorx.Wrap(e)' above
//
// | field optional | default value | example value
// - ${args_forbid_update} | "" | "user_id, game_id"
// - ${db_instance} "db.DB" | "db.DB"
// - ${handler_name} "HTTPListUser" | HTTPUpdateUser |
// - ${model} "model.User" | "payModel.Order"
// - ${handle_error} "fmt.Println(e, string(debug.Stack()))" | raven.Throw(e)
func GenerateUpdateOneAPI()
```
```go
package main
import (
    "encoding/json"
    "fmt"
    "github.com/fwhezfwhez/model_convert"
)
func main() {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int    `gorm:"column:state;default:" json:"state" form:"state"`
	}

	rs := GenerateUpdateOneAPI(VxTemplateUser{}, map[string]string{
		"${model}": "payModel.LingqianOrder",
		"${handler_name}" : "HTTPUpdateLingqianOrder",
		"${handle_error}": `common.SaveError(e)`,
		"${args_forbid_update}": "UserId, game_id",
	})
	fmt.Println(rs)
}
```
Output:
```go
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateUpdateOneAPI().
func HTTPUpdateLingqianOrder(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
        return
    }
    var count int
    if e:=db.DB.Model(&payModel.LingqianOrder{}).Where("id=?", idInt).Count(&count).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
        return
    }
    var param payModel.LingqianOrder
    if e:=c.Bind(&param);e!=nil {
        c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }

    if !util.IfZero(param.UserId) {
        c.JSON(400, gin.H{"message": "field 'UserId' can't be modified'"})
        return
    }
    if !util.IfZero(param.GameId) {
        c.JSON(400, gin.H{"message": "field 'GameId' can't be modified'"})
        return
    }
    if e:=db.DB.Model(&payModel.LingqianOrder{}).Where("id=?", id).Updates(param).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success"})
}
```

###### 2.5.6 generate crud

Note:
```go
// Generate update one api code.
// To completely use these code, you might import:
// "github.com/fwhezfwhez/errorx"
// you can get 'errorx.Wrap(e)' above
//
// | field optional | default value | example value
// - ${args_forbid_update} | "" | "user_id, game_id"
// - ${db_instance} "db.DB" | "db.DB"
// - ${handler_name} "HTTPListUser" | HTTPUpdateUser |
// - ${model} "model.User" | "payModel.Order"
// - ${handle_error} "fmt.Println(e, string(debug.Stack()))" | raven.Throw(e)
func GenerateCRUD()
```
```go
package main
import (
    "encoding/json"
    "fmt"
    "github.com/fwhezfwhez/model_convert"
)
func main() {
	type VxTemplateUser struct {
		Id     int    `gorm:"column:id;default:" json:"id" form:"id"`
		GameId int    `gorm:"column:game_id;default:" json:"game_id" form:"game_id"`
		UserId int    `gorm:"column:user_id;default:" json:"user_id" form:"user_id"`
		OpenId string `gorm:"column:open_id;default:" json:"open_id" form:"open_id"`

		TemplateId string `gorm:"column:template_id;default:" json:"template_id" form:"template_id"`
		State      int    `gorm:"column:state;default:" json:"state" form:"state"`
	}

	rs := GenerateCRUD(VxTemplateUser{}, map[string]string{
		"${model}": "moduleModel.VxTemplateUser",
		"${handle_error}": `common.SaveError(e)`,
		"${db_instance}": "db.DB",
	})
	fmt.Println(rs)
}
```
Output:
```go
// Auto generated by github.com/fwhezfwhez/model_convert.GenerateCRUD. You might need import:
// "github.com/gin-gonic/gin"
// "github.com/fwhezfwhez/errorx"
// "github.com/fwhezfwhez/model_convert/util"
//
// "package/path/to/db.DB"
// "package/path/to/moduleModel.VxTemplateUser"
//
// Auto generate by github.com/fwhezfwhez/model_convert.GenerateAddOneAPI().
func HTTPAddVxTemplateUser (c *gin.Context) {
    var param moduleModel.VxTemplateUser
    if e := c.Bind(&param); e!=nil {
        c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if e:=db.DB.Model(&moduleModel.VxTemplateUser{}).Create(&param).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success", "data": param})
}

// Auto generate by github.com/fwhezfwhez/model_convert.GenerateListAPI().
func HTTPListVxTemplateUser(c *gin.Context) {
    var engine = db.DB.Model(&moduleModel.VxTemplateUser{})

    id := c.DefaultQuery("id", "")
    if id != "" {
        engine = engine.Where("id = ?", id)
    }
    gameId := c.DefaultQuery("game_id", "")
    if gameId != "" {
        engine = engine.Where("game_id = ?", gameId)
    }
    userId := c.DefaultQuery("user_id", "")
    if userId != "" {
        engine = engine.Where("user_id = ?", userId)
    }
    openId := c.DefaultQuery("open_id", "")
    if openId != "" {
        engine = engine.Where("open_id = ?", openId)
    }
    templateId := c.DefaultQuery("template_id", "")
    if templateId != "" {
        engine = engine.Where("template_id = ?", templateId)
    }
    state := c.DefaultQuery("state", "")
    if state != "" {
        engine = engine.Where("state = ?", state)
    }
    page := c.DefaultQuery("page", "1")
    size := c.DefaultQuery("size", "20")
    orderBy := c.DefaultQuery("order_by", "")
    var count int
    if e:= engine.Count(&count).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    var list = make([]moduleModel.VxTemplateUser, 0, 20)
    if count == 0 {
        c.JSON(200, gin.H{"message": "success", "count": 0, "data": list})
        return
    }
    limit, offset := util.ToLimitOffset(size, page, count)
    engine = engine.Limit(limit).Offset(offset)
    if orderBy != "" {
        engine = engine.Order(util.GenerateOrderBy(orderBy))
    }
    if e:= engine.Find(&list).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success", "count": count, "data": list})
}

// Auto generate by github.com/fwhezfwhez/model_convert.GenerateGetOneAPI().
func HTTPGetVxTemplateUser(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
        return
    }
    var count int
    if e:=db.DB.Model(&moduleModel.VxTemplateUser{}).Where("id=?", idInt).Count(&count).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
        return
    }
    var instance moduleModel.VxTemplateUser
    if e:=db.DB.Model(&moduleModel.VxTemplateUser{}).Where("id=?", id).First(&instance).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success", "data": instance})
}

// Auto generate by github.com/fwhezfwhez/model_convert.GenerateUpdateOneAPI().
func HTTPUpdateVxTemplateUser(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
        return
    }
    var count int
    if e:=db.DB.Model(&moduleModel.VxTemplateUser{}).Where("id=?", idInt).Count(&count).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
        return
    }
    var param moduleModel.VxTemplateUser
    if e:=c.Bind(&param);e!=nil {
        c.JSON(400, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if e:=db.DB.Model(&moduleModel.VxTemplateUser{}).Where("id=?", id).Updates(param).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success"})
}

// Auto generate by github.com/fwhezfwhez/model_convert.GenerateDeleteOneAPI().
func HTTPDeleteVxTemplateUser(c *gin.Context) {
    id := c.Param("id")
    idInt, e := strconv.Atoi(id)
    if e!=nil {
        c.JSON(400, gin.H{"message": fmt.Sprintf("param 'id' requires int but got %s", id)})
        return
    }
    var count int
    if e:=db.DB.Model(&moduleModel.VxTemplateUser{}).Where("id=?", idInt).Count(&count).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    if count ==0 {
        c.JSON(200, gin.H{"message": fmt.Sprintf("id '%s' record not found", id)})
        return
    }
    var instance moduleModel.VxTemplateUser
    if e:=db.DB.Model(&moduleModel.VxTemplateUser{}).Where("id=?", id).Delete(&instance).Error; e!=nil {
        common.SaveError(e)
        c.JSON(500, gin.H{"message": errorx.Wrap(e).Error()})
        return
    }
    c.JSON(200, gin.H{"message": "success"})
}
```

#### 2.6 grpc

##### 2.6.1 grpc-proto-service to go model

```go
package main

import(
   ""
)
func main() {
	var src = `
service UserChargeLog{
    rpc GetUserChargeLog(ProMiniGameGetUserChargeLogRequest) returns (ProMiniGameGetUserChargeLogResponse) {}
}
`
	rs := model_convert.GenerateGRPCInstance(src, GenerateGRPCInstanceArg{
		PbPackagePath: "project/control/statistics/statisticsPb",
		PackageName:   "statisticsControl",
	})

	fmt.Println(rs)
}
```
Output:
```go
// auto generated by github.com/fwhezfwhez/model_convert.GenerateGRPCInstance
package statisticsControl

import (
	"context"
	// "github.com/fwhezfwhez/errorx"
	// "golang.org/x/protobuf/proto"
	"project/control/statistics/statisticsPb"
)

type UserChargeLog struct{}

func (o *UserChargeLog) GetUserChargeLog(ctx context.Context, param *statisticsPb.ProMiniGameGetUserChargeLogRequest) (*statisticsPb.ProMiniGameGetUserChargeLogResponse, error) {
	// todo do your work here
    rsp := &statisticsPb.ProMiniGameGetUserChargeLogResponse{
	}
	return rsp, nil
}

func (o *UserChargeLog) GetUserInfo(ctx context.Context, param *statisticsPb.ProMiniGameGetUserInfoRequest) (*statisticsPb.ProMiniGameGetUserInfoResponse, error) {
	// todo do your work here
    rsp := &statisticsPb.ProMiniGameGetUserInfoResponse{
	}
	return rsp, nil
}

func (o *UserChargeLog) GetUserCoin(ctx context.Context, param *statisticsPb.ProMiniGameGetUserCoinRequest) (*statisticsPb.ProMiniGameGetUserCoinResponse, error) {
	// todo do your work here
    rsp := &statisticsPb.ProMiniGameGetUserCoinResponse{
	}
	return rsp, nil
}
```