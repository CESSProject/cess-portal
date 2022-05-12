package command

import (
	"cess-portal/client"
	"github.com/spf13/cobra"
)

func NewQueryCommand() *cobra.Command {
	fc := &cobra.Command{
		Use:   "query <subcommand>",
		Short: "Query commands use for implement all of related find specific detail information",
	}

	fc.AddCommand(NewQueryPriceCommand())
	fc.AddCommand(NewQueryPurchasedSpaceCommand())
	fc.AddCommand(NewQueryFileCommand())

	return fc
}

// NewFindPriceCommand
func NewQueryPriceCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "price",
		Short: "Query the current storage space price(TCESS/GB)",
		Long:  `Price command chain to query and display the current storage space rental unit price (unit: MB).`,

		Run: QueryPriceCommandFunc,
	}

	return cc
}

func QueryPriceCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)

	client.QueryPrice()
}

func NewQueryPurchasedSpaceCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "space",
		Short: "Query real-time storage space detailed information",
		Long: `Space command chain query current account purchased storage space usage (used and remaining).
`,

		Run: QueryPurchasedSpaceCommand,
	}

	return cc
}

func QueryPurchasedSpaceCommand(cmd *cobra.Command, args []string) {
	InitComponents(cmd)

	client.QueryPurchasedSpace()
}

func NewQueryFileCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "file <fileId>",
		Short: "Query the uploaded files information",
		Long:  `File command chain query file information ,If <fileId> is null then show the file list`,

		Run: QueryFileCommand,
	}

	return cc
}

func QueryFileCommand(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	fileid := ""
	if len(args) != 0 {
		fileid = args[0]
	} else {
		cmd.Println("No parameter query, return a list of all files")
	}

	client.QueryFile(fileid)
}
