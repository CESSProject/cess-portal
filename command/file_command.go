package command

import (
	"cess-portal/client"
	"cess-portal/conf"
	"cess-portal/internal/logger"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewFileCommand() *cobra.Command {
	fc := &cobra.Command{
		Use:   "file <subcommand>",
		Short: "File commands use for implement related file function operate",
	}

	fc.AddCommand(NewFileUploadCommand())
	fc.AddCommand(NewFileDownloadCommand())
	fc.AddCommand(NewFileDeleteCommand())
	return fc
}

func NewFileUploadCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "upload <file path> <bucket name>",
		Short: "Upload the any specific file you want",
		Run:   FileUploadCommandFunc,
	}

	return cc
}

func FileUploadCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	if len(args) < 2 {
		fmt.Printf("Please enter correct parameters 'upload <file path> <bucket name>'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	client.FileUpload(args[0], args[1])
}

func NewFileDownloadCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "download <file id> <save directory>",
		Short: "Download the any specific file you want",
		Long:  `Download command mean download file from the CESS networks based on fileid, and save directory point where the downloaded file is saved.`,

		Run: FileDownloadCommandFunc,
	}

	return cc
}

func FileDownloadCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	if len(args) < 2 {
		fmt.Printf("Please enter the fileid and save directory of the download file 'file download <fileid> <save directory>'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	client.FileDownload(args[0], args[1])
}

func NewFileDeleteCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "delete <file id>",
		Short: "Delete the any specific file you want",
		Long:  `Delete command means removing the file from CESS networks`,

		Run: FileDeleteCommandFunc,
	}

	return cc
}

func FileDeleteCommandFunc(cmd *cobra.Command, args []string) {
	refreshProfile(cmd)
	logger.Log_Init()
	if len(args) == 0 {
		fmt.Printf("Please enter the fileid of the delete file'file delete <fileid>'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	client.FileDelete(args[0])
}
