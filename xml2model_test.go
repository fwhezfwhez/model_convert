package model_convert

import (
	"fmt"
	"testing"
)

func TestXMLToModel(t *testing.T) {
	fmt.Println(XMLToModel(
		`
<xml>
  <appid>wx977a002bb48c6a90</appid>
  <attach></attach>
  <bank_type>OTHERS</bank_type>
  <fee_type>CNY</fee_type>
  <is_subscribe>N</is_subscribe>
  <mch_id>1526815311</mch_id>
  <nonce_str>apmpnGV8MBFQ3wIdB0ELa7DdubvbjHfI</nonce_str>
  <openid>odOqR5kVI-z1D2xYzhDUJZWjQOrY</openid>
  <out_trade_no>1608011847_33873189_9526</out_trade_no>
  <result_code>SUCCESS</result_code>
  <err_code></err_code>
  <err_code_des></err_code_des>
  <return_code>SUCCESS</return_code>
  <return_msg></return_msg>
  <sign>E53660B71468A6A3B2497EAF04AD1472</sign>
  <time_end>20201215135733</time_end>
  <total_fee>10</total_fee>
  <coupon_fee></coupon_fee>
  <coupon_count></coupon_count>
  <coupon_type></coupon_type>
  <coupon_id></coupon_id>
  <trade_type>APP</trade_type>
  <transaction_id>4200000793202012151699848322</transaction_id>
</xml>
        `,
		"NotifyRequest",
	))
}
