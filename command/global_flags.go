package command

import (
	"cess-portal/conf"
	"cess-portal/internal/chain"
	"cess-portal/tools"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GlobalFlags struct {
	ConfFilePath string
}

func refreshProfile(cmd *cobra.Command) {
	configpath1, _ := cmd.Flags().GetString("config")
	configpath2, _ := cmd.Flags().GetString("c")
	if configpath1 != "" {
		conf.ConfigFilePath = configpath1
	} else {
		conf.ConfigFilePath = configpath2
	}
	parseProfile()
}

func parseProfile() {
	var (
		err          error
		confFilePath string
	)
	if conf.ConfigFilePath == "" {
		confFilePath = "./conf.toml"
	} else {
		confFilePath = conf.ConfigFilePath
	}
	f, err := os.Stat(confFilePath)
	if err != nil {
		log.Printf("[err] The '%v' file does not exist.\n", confFilePath)
		os.Exit(1)
	}
	if f.IsDir() {
		log.Printf("[err] The '%v' is not a file.\n", confFilePath)
		os.Exit(1)
	}

	viper.SetConfigFile(confFilePath)
	viper.SetConfigType("toml")

	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("[err] The '%v' file type error.\n", confFilePath)
		os.Exit(1)
	}

	err = viper.Unmarshal(conf.C)
	if err != nil {
		log.Printf("[err] Configuration file error, please use the default command to generate a template.\n")
		os.Exit(1)
	}

	if conf.C.RpcAddr == "" || conf.C.AccountSeed == "" || conf.C.AccountId == "" {
		log.Printf("[err] The configuration file cannot have empty entries.\n")
		os.Exit(1)
	}
	//
	if err := tools.CreatDirIfNotExist(conf.BaseDir); err != nil {
		log.Printf("[err] %v\n", err)
		os.Exit(1)
	}

	if err := tools.CreatDirIfNotExist(conf.LogfileDir); err != nil {
		log.Printf("[err] %v\n", err)
		os.Exit(1)
	}

	if err := tools.CreatDirIfNotExist(conf.FileCacheDir); err != nil {
		log.Printf("[err] %v\n", err)
		os.Exit(1)
	}
	//
	chain.ChainClient, err = chain.NewChainClient(conf.C.RpcAddr, conf.C.AccountSeed, conf.TimeToWaitEvents)
	if err != nil {
		log.Printf("[err] %v\n", err)
		os.Exit(1)
	}
	conf.PublicKey, err = tools.DecodePublicKeyOfCessAccount(conf.C.AccountId)
	if err != nil {
		log.Printf("[err] %v\n", err)
		os.Exit(1)
	}
	conf.PublicKeyfile = string(chain.ChainClient.GetPublicKey())
	if err != nil {
		log.Printf("[err] %v\n", err)
		os.Exit(1)
	}
}
