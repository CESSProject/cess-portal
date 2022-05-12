package command

import (
	"cess-portal/client"
	"cess-portal/conf"
	"cess-portal/tools"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"os"
)

func NewFileCommand() *cobra.Command {
	fc := &cobra.Command{
		Use:   "file <subcommand>",
		Short: "File commands use for implement related file function operate",
	}

	fc.AddCommand(NewFileUploadCommand())
	fc.AddCommand(NewFileDownloadCommand())
	fc.AddCommand(NewFileDeleteCommand())
	fc.AddCommand(NewFileDecodeCommand())

	return fc
}

func NewFileUploadCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "upload <filepath> <backups> <private key>",
		Short: "Upload the any specific file you want",
		Long:  `Upload command send local source files to scheduling nodes.`,

		Run: FileUploadCommandFunc,
	}

	return cc
}

func FileUploadCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	if len(args) < 2 {
		fmt.Printf("Please enter correct parameters 'upload <filepath> <backups> <private key>'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	fmt.Printf("%s[Warming] Do you want to upload your file without private key (it's means your file status is public)?%s\n", tools.Red, tools.Reset)
	fmt.Printf("%sYou can type the 'private key' or enter with nothing to skip it:%s", tools.Red, tools.Reset)
	psw, _ := gopass.GetPasswdMasked()

	client.FileUpload(args[0], args[1], string(psw))
}

func NewFileDownloadCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "download <fileid>",
		Short: "Download the any specific file you want",
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

func NewFileDeleteCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "delete <fileid>",
		Short: "Delete the any specific file you want",
		Long:  `Deleting a file means removing the file from CESS,But there may be a delay.`,

		Run: FileDeleteCommandFunc,
	}

	return cc
}

func FileDeleteCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	if len(args) == 0 {
		fmt.Printf("Please enter the fileid of the deleted file'delete <fileid>'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	client.FileDelete(args[0])

}

func NewFileDecodeCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "decode <filepath>",
		Short: "Decrypt the any specific file again when you failed file decrypt first chance",
		Long:  `File decode means that if the file is not decrypted when you download it, it can be decode by this method.Please enter absolute path.`,

		Run: FileDecodeCommandFunc,
	}

	return cc
}

func FileDecodeCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)

	if len(args) == 0 {
		fmt.Printf("Please enter the path of the file to be decoded'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	client.FileDecode(args[0])
}
