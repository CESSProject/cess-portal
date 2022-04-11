package command

import (
	"cess-portal/client"
	"github.com/spf13/cobra"
)

func NewFindCommand() *cobra.Command {
	fc := &cobra.Command{
		Use:   "find <subcommand>",
		Short: "Storage related commands",
	}

	fc.AddCommand(NewFindPriceCommand())
	fc.AddCommand(NewFindPurchasedSpaceCommand())
	fc.AddCommand(NewFindFileCommand())

	return fc
}

// NewFindPriceCommand
func NewFindPriceCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "price",
		Short: "Price refers to the query storage price",
		Long:  `Price command chain to query and display the current storage space rental unit price (unit: MB).`,

		Run: FindPriceCommandFunc,
	}

	return cc
}

func FindPriceCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)

	client.FindPrice()
}

func NewFindPurchasedSpaceCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "space",
		Short: "Space refers to the query storage space",
		Long: `Space command chain query current account purchased storage space usage (used and remaining).
`,

		Run: FindPurchasedSpaceCommand,
	}

	return cc
}

func FindPurchasedSpaceCommand(cmd *cobra.Command, args []string) {
	InitComponents(cmd)

	client.FindPurchasedSpace()
}

func NewFindFileCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "file <fileId>",
		Short: "File refers to the query storage file",
		Long:  `File command chain query file information ,If <fileId> is null then show the file list`,

		Run: FindFileCommand,
	}

	return cc
}

func FindFileCommand(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	fileid := ""
	if len(args) != 0 {
		fileid = args[0]
	} else {
		cmd.Println("No parameter query, return a list of all files")
	}

	client.FindFile(fileid)
}
