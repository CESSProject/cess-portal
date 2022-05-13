package command

import (
	"cess-portal/client"
	"cess-portal/conf"
	"cess-portal/internal/logger"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

func NewPurchaseCommand() *cobra.Command {
	tc := &cobra.Command{
		Use:   "trade <subcommand>",
		Short: "Purchase commands use for implement all of related transaction function",
	}

	tc.AddCommand(NewPurchaseBuySpaceCommand())
	tc.AddCommand(NewPurchaseObtainCommand())

	return tc
}

func NewPurchaseBuySpaceCommand() *cobra.Command {
	tbs := &cobra.Command{
		Use:   "storage <spacequantity>  <expected price>",
		Short: "Buy CESS storage space",
		Long:  `<spacequantity> storage space quantity you want to buy，unit(TCESS/GB); <space duration> storage space you want to rental, unit(TCESS/Month); <expected price> set the expected price(integer),unit(TCESS/GB) for the purchase, if input null mean accept the CESS real-time storage unit price.`,

		Run: PurchaseBuySpaceCommandFunc,
	}

	return tbs
}

func PurchaseBuySpaceCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	var expected = 0
	var quantity = 0
	var duration = 0
	var err error
	if len(args) < 2 {
		fmt.Printf("[Error]Please fill in the amount of storage space you want to purchase! Usage: cessctl trade exp <quantity> <duration>\n")
		logger.OutPutLogger.Sugar().Infof("[Error]Please fill in the amount of storage space you want to purchase! Usage: cessctl trade exp <spacequantity> <duration>\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	if len(args) > 2 {
		expected, err = strconv.Atoi(args[2])
		if err != nil || expected < 0 {
			fmt.Printf("[Error]Please enter the correct number (integer) in <expected price>\n")
			logger.OutPutLogger.Sugar().Infof("[Error]Please enter the correct number (integer) in <expected price>\n")
			os.Exit(conf.Exit_CmdLineParaErr)
		}
	}
	quantity, err1 := strconv.Atoi(args[0])
	duration, err2 := strconv.Atoi(args[1])
	if err1 != nil || err2 != nil || quantity < 0 {
		fmt.Printf("[Error]Please enter the correct number (integer) in <spacequantity>\n")
		logger.OutPutLogger.Sugar().Infof("[Error]Please enter the correct number (integer) in <spacequantity>\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	client.Expansion(quantity, duration, expected)
}

func NewPurchaseObtainCommand() *cobra.Command {
	tbs := &cobra.Command{
		Use:   "free <address>",
		Short: "Top up free TCESS from the faucet",
		Long:  `Free command use for obtain the TCESS tokens from the CESS faucet service, the amount of each top up is 10000 TCESS`,

		Run: PurchaseObtainCommandFunc,
	}

	return tbs
}

func PurchaseObtainCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	if len(args) == 0 {
		fmt.Printf("[Error]Please fill in the account public key! Usage: cessctl trade obtain <public key>")
		logger.OutPutLogger.Sugar().Infof("Please fill in the account public key! Usage: cessctl trade obtain <public key>\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	if len(args[0]) != 66 {
		fmt.Printf("[Error]Please enter the correct number of digits for the public key!\nThe way to get public key——>>https://polkadot.subscan.io/tools/ss58_transform")
		logger.OutPutLogger.Sugar().Infof("[Error]Please enter the correct number of digits for the public key!\nThe way to get public key——>>https://polkadot.subscan.io/tools/ss58_transform\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	if !strings.HasPrefix(args[0], "0x") {
		fmt.Println("[Error]The public key you entered is not in the correct format!\nThe way to get public key——>>https://polkadot.subscan.io/tools/ss58_transform")
		logger.OutPutLogger.Sugar().Infof("[Error]The public key you entered is not in the correct format!\nThe way to get public key——>>https://polkadot.subscan.io/tools/ss58_transform\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	client.ObtainFromFaucet(args[0])
}
