package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestFileUpload(t *testing.T) {
	conf.ClientConf.ChainData.CessRpcAddr = ""
	conf.ClientConf.BoardInfo.BoardPath = ""
	path := ""
	backups := ""
	PrivateKey := ""

	conf.ClientConf.PathInfo.KeyPath = ""
	conf.ClientConf.ChainData.IdAccountPhraseOrSeed = ""
	client.FileUpload(path, backups, PrivateKey)
}
