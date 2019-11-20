package model_convert

import (
	"fmt"
	"github.com/robfig/cron"
	"testing"
)

func TestLowerFistLetter(t *testing.T) {
    c := cron.New()
    c.AddFunc("* * * ? * 1", func() {
		fmt.Println(1)
	})
    c.Start()

	select {

	}
}
