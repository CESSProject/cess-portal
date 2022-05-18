package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestObtainFromFaucet(t *testing.T) {
	//config file
	conf.ClientConf.BoardInfo.BoardPath = ""
	conf.ClientConf.ChainData.FaucetAddress = "http://139.224.19.104:9708/transfer"

	//param
	pkg := "cXhqoREguFc6SteocXoiyhjoeJXzF4yng71wd2fZn2jVrzxSy"
	err := client.ObtainFromFaucet(pkg)
	t.Fatal(err)
}
