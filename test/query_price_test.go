package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestFindPrice(t *testing.T) {
	//config file
	conf.ClientConf.ChainData.CessRpcAddr = ""
	conf.ClientConf.ChainData.AccountPublicKey = ""
	conf.ClientConf.BoardInfo.BoardPath = ""

	err := client.QueryPrice()
	t.Fatal(err)
}
