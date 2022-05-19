package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestDownload(t *testing.T) {
	//config file
	conf.ClientConf.ChainData.CessRpcAddr = ""
	conf.ClientConf.BoardInfo.BoardPath = ""
	conf.ClientConf.PathInfo.InstallPath = ""

	//param
	fileid := ""

	err := client.FileDownload(fileid)
	t.Fatal(err)
}
