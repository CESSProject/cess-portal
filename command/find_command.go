package command

import (
	"fmt"
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
		Long: `Price command chain to query and display the current storage space rental unit price (unit: MB).
`,

		Run: FindPriceCommandFunc,
	}

	return cc
}

func FindPriceCommandFunc(cmd *cobra.Command, args []string) {
	cmd.Printf("there is find price command!\n")
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
	fmt.Println("there is Purchased Space command!")
}

func NewFindFileCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "file",
		Short: "File refers to the query storage file",
		Long:  `File command chain query all file information that has been uploaded by the current account (sorting, keyword retrieval...).`,

		Run: FindFindFileCommand,
	}

	return cc
}

func FindFindFileCommand(cmd *cobra.Command, args []string) {
	fmt.Println("there is Find File command!")
}
