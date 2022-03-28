package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestFindFile(t *testing.T) {
	//config file
	conf.ClientConf.ChainData.CessRpcAddr = ""
	conf.ClientConf.ChainData.AccountPublicKey = ""
	conf.ClientConf.BoardInfo.BoardPath = ""

	//param
	fileid := ""

	err := client.FindFile(fileid)
	t.Fatal(err)
}
