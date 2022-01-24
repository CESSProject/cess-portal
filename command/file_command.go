package command

import (
	"dapp_cess_client/client"
	"dapp_cess_client/conf"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func NewFileCommand() *cobra.Command {
	fc := &cobra.Command{
		Use:   "file <subcommand>",
		Short: "File related commands",
	}

	fc.AddCommand(NewFileUploadCommand())
	fc.AddCommand(NewFileDownloadCommand())

	return fc
}

func NewFileUploadCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "upload <filepath> <downloadfee>",
		Short: "upload refers to the upload file",
		Long:  `Price command send local source files to scheduling nodes.`,

		Run: FileUploadCommandFunc,
	}

	return cc
}

func FileUploadCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	if len(args) < 2 {
		fmt.Printf("Please enter correct parameters 'upload <filepath> <downloadfee>'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	client.FileUpload(args[0], args[1])
}

func NewFileDownloadCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "download <fileid>",
		Short: "download refers to the download file",
		Long:  `Download command download file based on fileId.`,

		Run: FileDownloadCommandFunc,
	}

	return cc
}

func FileDownloadCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	if len(args) == 0 {
		fmt.Printf("Please enter the fileid of the downloaded file 'download <fileid>'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	client.FileDownload(args[0])
}
