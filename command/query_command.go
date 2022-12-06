package command

import (
	"cess-portal/client"
	"cess-portal/conf"
	"cess-portal/internal/logger"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewQueryCommand() *cobra.Command {
	fc := &cobra.Command{
		Use:   "query <subcommand>",
		Short: "Query commands use for implement all of related find specific detail information",
	}

	fc.AddCommand(NewQueryFilestateCommand())
	fc.AddCommand(NewQueryFilelistCommand())
	fc.AddCommand(NewQueryBucketlistCommand())
	fc.AddCommand(NewQuerySpaceCommand())
	return fc
}

func NewQueryFilestateCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "fstate <file id>",
		Short: "Query the state of files in the CESS system",
		Run:   QueryFilestateCommandFunc,
	}

	return cc
}
func NewQueryFilelistCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "files <bucket name>",
		Short: "Query the list of files in the CESS system",
		Run:   QueryFilelistCommandFunc,
	}

	return cc
}

func NewQueryBucketlistCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "buckets",
		Short: "Query the list of bucket in the CESS system",
		Run:   QueryBucketlistCommandFunc,
	}

	return cc
}

func NewQuerySpaceCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "space",
		Short: "Query space info of you account in the CESS system",
		Run:   QuerySpaceCommandFunc,
	}
	return cc
}

func QuerySpaceCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	client.UserSpaceQuery()
}

func QueryFilestateCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	if len(args) < 1 {
		fmt.Printf("Please enter the file id.\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	client.FilestateQuery(args[0])
}

func QueryFilelistCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	if len(args) < 1 {
		fmt.Printf("Please enter the bucket name.\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	client.FilelistQuery(args[0])
}

func QueryBucketlistCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	client.BucketlistQuery()
}
