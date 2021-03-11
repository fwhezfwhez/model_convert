package model_convert

import (
	"fmt"
	"testing"
)

func TestProf(t *testing.T) {
	ngconf, gosdk := GenerateProf("10.10.12.12", ":6060", "ddz.prof.com")

	fmt.Println(ngconf)

	fmt.Println(gosdk)
}
