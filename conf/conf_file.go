package conf

import (
	"cess-portal/tools"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type CessClient struct {
	BoardInfo BoardInfo `yaml:"boardInfo"`
	ChainData ChainData `yaml:"chainData"`
	PathInfo  PathInfo  `yaml:"pathInfo"`
}

type BoardInfo struct {
	BoardPath string `yaml:"boardPath"`
}
type ChainData struct {
	CessRpcAddr           string `yaml:"cessRpcAddr"`
	FaucetAddress         string `yaml:"faucetAddress"`
	IdAccountPhraseOrSeed string `yaml:"idAccountPhraseOrSeed"`
	AccountPublicKey      string `yaml:"accountPublicKey"`
	WalletAddress         string `yaml:"walletAddress"`
}
type PathInfo struct {
	KeyPath     string `yaml:"keyPath"`
	InstallPath string `yaml:"installPath"`
}

var ClientConf = new(CessClient)
var ConfFilePath string

func InitConf() {
	if ConfFilePath == "" {
		ConfFilePath = Conf_File_Path_D
	}
	_, err := os.Stat(ConfFilePath)
	if err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m The '%v' config file does not exist\n", 41, ConfFilePath)
		os.Exit(Exit_CmdLineParaErr)
	}
	yamlFile, err := ioutil.ReadFile(ConfFilePath)
	if err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m The '%v' file read error\n", 41, ConfFilePath)
		os.Exit(Exit_ConfErr)
	}
	err = yaml.Unmarshal(yamlFile, ClientConf)
	if err != nil {
		fmt.Printf("\x1b[%dm[err]\x1b[0m The '%v' file format error\n", 41, ConfFilePath)
		os.Exit(Exit_ConfErr)
	}
	pubkey, err := tools.DecodeToPub(ClientConf.ChainData.WalletAddress, tools.ChainCessTestPrefix)
	if err != nil {
		fmt.Printf("[Error]The wallet address you entered is incorrect, please re-enter\n")
		os.Exit(Exit_ConfErr)
	}
	ClientConf.ChainData.AccountPublicKey = tools.PubBytesTo0XString(pubkey)
}
