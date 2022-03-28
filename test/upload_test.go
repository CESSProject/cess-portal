package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestFileUpload(t *testing.T) {
	//config file
	conf.ClientConf.ChainData.CessRpcAddr = ""
	conf.ClientConf.BoardInfo.BoardPath = ""
	conf.ClientConf.PathInfo.KeyPath = ""
	conf.ClientConf.ChainData.IdAccountPhraseOrSeed = ""
	conf.ClientConf.ChainData.WalletAddress = ""

	//param
	path := ""
	backups := ""
	PrivateKey := ""
	client.FileUpload(path, backups, PrivateKey)
}
