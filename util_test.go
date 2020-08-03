package model_convert

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron"
)

func TestLowerFistLetter(t *testing.T) {
	c := cron.New()
	c.AddFunc("* * * ? * 1", func() {
		fmt.Println(1)
	})
	c.Start()

	select {
	case <-time.After(5 * time.Second):

	}
}

func TestRandomNumber(t *testing.T) {
	var m = map[string]interface{}{
		"name": "ft",
	}
	b, _ := json.MarshalIndent(m, "    ", "    ")
	fmt.Println(string(b))
}

func TestPicked(t *testing.T) {
	var count int
	for i := 0; i < 1000000; i++ {
		if Pick(0.000001) {
			count++
		}
	}
	fmt.Println(count)
}
