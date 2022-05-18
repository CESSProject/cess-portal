package test

import (
	"cess-portal/client"
	"cess-portal/conf"
	"testing"
)

func TestFileDecrypt(t *testing.T) {
	conf.ClientConf.PathInfo.InstallPath = ""
	decryptpath := ""
	err := client.FileDecrypt(decryptpath)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
