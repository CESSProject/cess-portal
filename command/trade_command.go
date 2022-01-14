package command

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewTradeCommand() *cobra.Command {
	tc := &cobra.Command{
		Use:   "trade <subcommand>",
		Short: "Trade related commands",
	}

	tc.AddCommand(NewTradeBuySpaceCommand())
	tc.AddCommand(NewFileDownloadCommand())

	return tc
}

func NewTradeBuySpaceCommand() *cobra.Command {
	tbs := &cobra.Command{
		Use:   "trade <spacesize>",
		Short: "trade refers to the trade with cess chian",
		Long:  `Trade command send on-chain transactions, buy space.`,

		Run: TradeBuySpaceCommandFunc,
	}

	return tbs
}

func TradeBuySpaceCommandFunc(cmd *cobra.Command, args []string) {
	fmt.Println("there is Trade Buy Space command!")
}
