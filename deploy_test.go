package model_convert

import (
	"fmt"
	"testing"
)

func TestGenerateDeployNginxConf(t *testing.T) {
	rs := GenerateDeployNginxConf("xyx.zonst.com", ":11115", "xyx_prop_present", "/data/xxxx.log", "/data/nginx-alives")

	fmt.Println(rs)
}

func TestGenerateDeploySupervisorIni(t *testing.T) {
	rs :=GenerateDeploySupervisorIni("xyx_prop_present", "/bra/bra", "-mode 'pro' -p ':8080'", "/data/log/bra/bra.log")
	fmt.Println(rs)
}
