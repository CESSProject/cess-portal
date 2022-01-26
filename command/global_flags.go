package command

import (
	"dapp_cess_client/conf"
	"dapp_cess_client/internal/logger"
	"github.com/spf13/cobra"
)

type GlobalFlags struct {
	ConfFilePath string
}

func InitComponents(cmd *cobra.Command) {
	configpath1, _ := cmd.Flags().GetString("config")
	configpath2, _ := cmd.Flags().GetString("c")
	if configpath1 != "" {
		conf.ConfFilePath = configpath1
	} else {
		conf.ConfFilePath = configpath2
	}
	conf.InitConf()
	logger.InitLogger()
}