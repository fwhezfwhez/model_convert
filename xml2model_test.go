package model_convert

import (
	"fmt"
	"testing"
)

func TestXMLToModel(t *testing.T) {
	fmt.Println(233/100*100)
	fmt.Println(XMLToModel(
		`
<xml> // 按照格式补充 

<return_code><![CDATA[SUCCESS]]></return_code> 

<return_msg><![CDATA[获取成功]]></return_msg> 

<result_code><![CDATA[SUCCESS]]></result_code> 

<mch_id>10000098</mch_id> 

<appid><![CDATA[wxe062425f740c30d8]]></appid> 

<detail_id><![CDATA[1000000000201503283103439304]]></detail_id> 

<partner_trade_no><![CDATA[1000005901201407261446939628]]></partner_trade_no> 

<status><![CDATA[SUCCESS]]></status>

<payment_amount>650</payment_amount > 

<openid ><![CDATA[oxTWIuGaIt6gTKsQRLau2M0yL16E]]></openid> 

<transfer_time><![CDATA[2015-04-21 20:00:00]]></transfer_time>

<transfer_name ><![CDATA[测试]]></transfer_name > 

<desc><![CDATA[福利测试]]></desc> 

</xml>

        `,
		"VXResponse",
	))
}
