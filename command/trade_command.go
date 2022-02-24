package command

import (
	"c-portal/client"
	"c-portal/conf"
	"c-portal/internal/logger"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
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
		Use:   "exp <spacequantity> <expected price>",
		Short: "exp refers to make your space bigger,unit:[1/512G].",
		Long:  `<spacequantity> buy space(not nullable); <expected price> set the expected price(integer) for the purchase(nullable) if null accept the storage unit price at the current moment.`,

		Run: TradeBuySpaceCommandFunc,
	}

	return tbs
}

func TradeBuySpaceCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	var expected = 0
	var quantity = 0
	var err error
	if len(args) == 0 {
		fmt.Printf("[Error]Please fill in the amount of storage space you want to purchase! Usage: cessctl trade exp <quantity>\n")
		logger.OutPutLogger.Sugar().Infof("[Error]Please fill in the amount of storage space you want to purchase! Usage: cessctl trade exp <spacequantity>\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	if len(args) > 1 {
		expected, err = strconv.Atoi(args[1])
		if err != nil || expected < 0 {
			fmt.Printf("[Error]Please enter the correct number (integer) in <expected price>\n")
			logger.OutPutLogger.Sugar().Infof("[Error]Please enter the correct number (integer) in <expected price>\n")
			os.Exit(conf.Exit_CmdLineParaErr)
		}
	}
	for _, r := range args[0] {
		if !unicode.IsNumber(r) {
			fmt.Printf("[Error]Please enter the number in <spacequantity>!\n")
			logger.OutPutLogger.Sugar().Infof("[Error]Please enter the number in <spacequantity>!\n")
			os.Exit(conf.Exit_CmdLineParaErr)
		}
	}
	quantity, err = strconv.Atoi(args[0])
	if err != nil || quantity < 0 {
		fmt.Printf("[Error]Please enter the correct number (integer) in <spacequantity>\n")
		logger.OutPutLogger.Sugar().Infof("[Error]Please enter the correct number (integer) in <spacequantity>\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	client.Expansion(quantity, expected)
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
