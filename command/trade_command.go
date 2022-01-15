package command

import (
	"dapp_cess_client/client"
	"dapp_cess_client/conf"
	"dapp_cess_client/internal/logger"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"unicode"
)

func NewTradeCommand() *cobra.Command {
	tc := &cobra.Command{
		Use:   "trade <subcommand>",
		Short: "Trade related commands",
	}

	tc.AddCommand(NewTradeBuySpaceCommand())
	tc.AddCommand(NewTradeObtainCommand())

	return tc
}

func NewTradeBuySpaceCommand() *cobra.Command {
	tbs := &cobra.Command{
		Use:   "exp <spacequantity>",
		Short: "exp refers to make your space bigger unit:[1/512G]",
		Long:  `Exp command send on-chain transactions, buy space.`,

		Run: TradeBuySpaceCommandFunc,
	}

	return tbs
}

func TradeBuySpaceCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	if len(args) == 0 {
		fmt.Printf("[Error]Please fill in the amount of storage space you want to purchase! Usage: cessctl trade exp <quantity>\n")
		logger.OutPutLogger.Sugar().Infof("[Error]Please fill in the amount of storage space you want to purchase! Usage: cessctl trade exp <quantity>\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	for _, r := range args[0] {
		if !unicode.IsNumber(r) {
			fmt.Printf("[Error]Please enter the number!\n")
			logger.OutPutLogger.Sugar().Infof("[Error]Please enter the number!\n")
			os.Exit(conf.Exit_CmdLineParaErr)
		}
	}
	client.Expansion(args[0])
}

func NewTradeObtainCommand() *cobra.Command {
	tbs := &cobra.Command{
		Use:   "obtain <address>",
		Short: "obtain refers to the trade with cess chian",
		Long:  `Obtain command get a certain amount of tokens through the faucet service.`,

		Run: TradeObtainCommandFunc,
	}

	return tbs
}

func TradeObtainCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	if len(args) == 0 {
		fmt.Printf("[Error]Please fill in the account public key! Usage: cessctl trade obtain <public key>\n")
		logger.OutPutLogger.Sugar().Infof("Please fill in the account public key! Usage: cessctl trade obtain <public key>\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	if len(args[0]) != 66 {
		fmt.Printf("[Error]Please enter the correct number of digits for the public key!\nThe way to get public key——>>https://polkadot.subscan.io/tools/ss58_transform\n")
		logger.OutPutLogger.Sugar().Infof("[Error]Please enter the correct number of digits for the public key!\nThe way to get public key——>>https://polkadot.subscan.io/tools/ss58_transform\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	if !strings.HasPrefix(args[0], "0x") {
		fmt.Println("[Error]The public key you entered is not in the correct format!\nThe way to get public key——>>https://polkadot.subscan.io/tools/ss58_transform\n")
		logger.OutPutLogger.Sugar().Infof("[Error]The public key you entered is not in the correct format!\nThe way to get public key——>>https://polkadot.subscan.io/tools/ss58_transform\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	client.ObtainFromFaucet(args[0])
}
