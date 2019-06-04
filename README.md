## model-convert

model-convert is used for transfer all kinds of structs

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Declaration](#declaration)
- [1. Start](#1-start)
- [2. Stable cases](#2-stable-cases)
    - [2.1 xml to go model](#21-xml-to-go-model)
    - [2.2 Gorm-postgres table to model](#22-gorm-postgres-table-to-model)
    - [2.3 Add json,form tag for go model](#23-add-jsonform-tag-for-go-model)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## Declaration
** Most cases are supposed to use before app starts, they're not advised using in runtime.It helps auto-build structures.
** Most structure auto-built might need to be formatted yourself.

## 1. Start
`go get github.com/fwhezfwhez/model_convert`

## 2. Stable cases
Most cases are developing.However, only when specific requirements are met with, I will upgrade requiring functions. Here are stable use cases, it will be taken care of when project are updating.

#### 2.1 xml to go model
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

#### 2.2 Gorm-postgres table to model

```go
package main

import (
    "fmt"
    "model_convert"
)

func main() {
    dataSouce := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", "localhost", "5432", "postgres", "game", "disable", "123")
    tableName := "user_info"
    fmt.Println(model_convert.TableToStructWithTag(dataSouce, tableName))
}
```
output:
```
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
