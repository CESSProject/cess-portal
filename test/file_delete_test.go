package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestFileDelete(t *testing.T) {
	//config file
	conf.ClientConf.ChainData.CessRpcAddr = ""
	conf.ClientConf.ChainData.IdAccountPhraseOrSeed = ""
	conf.ClientConf.BoardInfo.BoardPath = ""

	//param
	fileid := ""

	err := client.FileDelete(fileid)
	t.Fatal(err)
}
