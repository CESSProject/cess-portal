package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestFindFile(t *testing.T) {
	//config file
	conf.ClientConf.ChainData.CessRpcAddr = ""
	conf.ClientConf.BoardInfo.BoardPath = ""
	conf.ClientConf.ChainData.AccountPublicKey = ""

	//param
	fileid := ""

	err := client.QueryFile(fileid)
	t.Fatal(err)
}
