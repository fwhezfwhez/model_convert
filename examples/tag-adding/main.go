package main

import (
	"fmt"
	"model_convert"
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
