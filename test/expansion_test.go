package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestExpansion(t *testing.T) {
	//config file
	conf.ClientConf.ChainData.CessRpcAddr = ""
	conf.ClientConf.ChainData.IdAccountPhraseOrSeed = ""
	conf.ClientConf.BoardInfo.BoardPath = ""

	//param
	quantity := 1
	duration := 1
	expected := 0
	err := client.Expansion(quantity, duration, expected)
	t.Fatal(err)
}
