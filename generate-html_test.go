package model_convert

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/russross/blackfriday"
)

func TestGenHTML(t *testing.T) {

	buf, e := ioutil.ReadFile(`G:/gitbooks/Library/Import/xyx_games_doc/zhong-zhi-shang-rao-ma-jiang/httpwen-dang-shuo-ming/dui-jie/lian-chu-bei.md`)
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	output := blackfriday.MarkdownCommon(buf)

	fmt.Println(string(output))

	if e:= ioutil.WriteFile("./tmp/tmp.html", output, 0666); e!=nil {
		fmt.Println(e.Error())
		return
	}
}
