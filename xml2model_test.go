package model_convert

import (
	"fmt"
	"testing"
)

func TestXMLToModel(t *testing.T) {
	fmt.Println(233/100*100)
	fmt.Println(XMLToModel(
		`
            <xml>
                <ToUserName><![CDATA[公众号]]></ToUserName>
                <FromUserName><![CDATA[粉丝号]]></FromUserName>
                <CreateTime>1460537339</CreateTime>
                <MsgType><![CDATA[text]]></MsgType>
                <Content><![CDATA[欢迎开启公众号开发者模式]]></Content>
                <MsgId>6272960105994287618</MsgId>
            </xml>
        `,
		"TextMessage",
	))
}
