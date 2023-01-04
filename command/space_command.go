package command

import (
	"cess-portal/client"
	"cess-portal/conf"
	"cess-portal/internal/logger"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func NewSpaceCommand() *cobra.Command {
	tc := &cobra.Command{
		Use:   "space <subcommand>",
		Short: "space commands use for implement all of related transaction function",
	}
	tc.AddCommand(
		NewPurchaseSpaceCommand(),
		NewAuthSpaceCommand(),
		NewCancelAuthCommand(),
	)
	return tc
}

func NewPurchaseSpaceCommand() *cobra.Command {
	tbs := &cobra.Command{
		Use:   "purchase <space quantity>",
		Short: "purchase CESS storage space",
		Long:  `<space quantity> storage space quantity you want to purchase,unit:GiB`,
		Run:   PurchaseSpaceCommandFunc,
	}

	return tbs
}

func PurchaseSpaceCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	if len(args) > 0 {
		size, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			fmt.Println("Illegal page size")
			os.Exit(conf.Exit_CmdLineParaErr)
		}
		client.StoragePurchase(uint32(size))
	}
	fmt.Println("Illegal space size")
	os.Exit(conf.Exit_CmdLineParaErr)
}

func NewAuthSpaceCommand() *cobra.Command {
	tbs := &cobra.Command{
		Use:   "auth",
		Short: "authorize CESS storage space",
		Run:   AuthSpaceCommandFunc,
	}
	return tbs
}

func AuthSpaceCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	client.SpaceAuthorize()
}

func NewCancelAuthCommand() *cobra.Command {
	tbs := &cobra.Command{
		Use:   "cancel",
		Short: "cancel authorizition CESS storage space",
		Run:   AuthSpaceCommandFunc,
	}
	return tbs
}

func CancelAuthCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	client.AuthCancel()
}
