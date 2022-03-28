package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestObtainFromFaucet(t *testing.T) {
	//config file
	conf.ClientConf.BoardInfo.BoardPath = ""
	conf.ClientConf.ChainData.FaucetAddress = ""

	//param
	pkg := ""
	err := client.ObtainFromFaucet(pkg)
	t.Fatal(err)
}
