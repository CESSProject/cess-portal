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
		Use:   "purchase <subcommand>",
		Short: "Purchase commands use for implement all of related transaction function",
	}

	tc.AddCommand(NewPurchaseBuySpaceCommand())
	tc.AddCommand(NewPurchaseObtainCommand())

	return tc
}

func NewPurchaseBuySpaceCommand() *cobra.Command {
	tbs := &cobra.Command{
		Use:   "storage <space quantity> <space duration> <expected price>",
		Short: "Buy CESS storage space",
		Long:  `<space quantity> storage space quantity you want to buy,unit:GB; <space duration> storage space you want to rental, unit:Month; <expected price> set the expected price(integer),unit(TCESS/GB) for the purchase, if input null mean accept the CESS real-time storage unit price.`,

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
		fmt.Printf("[Error]Please fill in the amount of storage space you want to purchase! Usage: cessctl purchase storage <space quantity> <space duration>\n")
		logger.OutPutLogger.Sugar().Infof("[Error]Please fill in the amount of storage space you want to purchase! Usage: cessctl purchase storage <space quantity> <space duration>\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	if len(args) > 2 {
		expected, err = strconv.Atoi(args[2])
		if err != nil || expected < 0 {
			fmt.Printf("[Error]Please enter the correct number (integer) in <expected price> <space duration>\n")
			logger.OutPutLogger.Sugar().Infof("[Error]Please enter the correct number (integer) in <expected price> or <space duration>\n")
			os.Exit(conf.Exit_CmdLineParaErr)
		}
	}
	quantity, err1 := strconv.Atoi(args[0])
	duration, err2 := strconv.Atoi(args[1])
	if err1 != nil || err2 != nil || quantity < 0 {
		fmt.Printf("[Error]Please enter the correct number (integer) in <space quantity> or \n")
		logger.OutPutLogger.Sugar().Infof("[Error]Please enter the correct number (integer) in <space quantity>\n")
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
		fmt.Printf("[Error]Please fill in the account public key! Usage: cessctl trade obtain <wallet address>")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	if len(args[0]) != 49 {
		fmt.Printf("[Error]Please enter the correct number of digits for the wallet address!\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	if !strings.HasPrefix(args[0], "c") {
		fmt.Println("[Error]The wallet address you entered is not in the correct format!\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	client.ObtainFromFaucet(args[0])
}
