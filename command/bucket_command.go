package command

import (
	"cess-portal/client"
	"cess-portal/conf"
	"cess-portal/internal/logger"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewBucketCommand() *cobra.Command {
	fc := &cobra.Command{
		Use:   "bucket <subcommand>",
		Short: "bucket commands use for implement all of related find specific detail information",
	}

	fc.AddCommand(NewBucketCreateCommand())
	fc.AddCommand(NewBucketDeleteCommand())
	return fc
}

func NewBucketCreateCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "create <bucket name>",
		Short: "create bucket in the CESS system",
		Run:   CreateBucketCommandFunc,
	}

	return cc
}
func NewBucketDeleteCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "delete <bucket name>",
		Short: "delete bucket from the CESS system",
		Run:   DeleteBucketCommandFunc,
	}

	return cc
}

func CreateBucketCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	if len(args) < 1 {
		fmt.Printf("Please enter the bucket name.\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	client.BucketCreate(args[0])
}

func DeleteBucketCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	if len(args) < 1 {
		fmt.Printf("Please enter the bucket name.\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	client.BucketDelete(args[0])
}
