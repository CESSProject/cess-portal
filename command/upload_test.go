package command

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestFileUpload(t *testing.T) {
	conf.ClientConf.ChainData.CessRpcAddr = ""
	path := ""
	backups := ""
	PrivateKey := ""

	conf.ClientConf.KeyInfo.KeyPath = ""
	conf.ClientConf.ChainData.IdAccountPhraseOrSeed = ""
	client.FileUpload(path, backups, PrivateKey)
}
