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
	fc.AddCommand(NewFileDecryptCommand())

	return fc
}

func NewFileUploadCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "upload <file path> <backups>",
		Short: "Upload the any specific file you want",
		Long:  `Upload command mean send the local source files to CESS nework scheduling nodes;You can input any 16/24/32 length numbers to be your private key, then others people unable to decrypt your file data. if you choose private key is nil, then system is default you file become public file.`,

		Run: FileUploadCommandFunc,
	}

	return cc
}

func FileUploadCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	if len(args) < 2 {
		fmt.Printf("Please enter correct parameters 'upload <filepath> <backups>'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	fmt.Printf("%s[Warming] Do you want to upload your file without private key (it's means your file status is public)?%s\n", tools.Red, tools.Reset)
	fmt.Printf("%sYou can type the 'private key' or enter with nothing to skip it:%s", tools.Red, tools.Reset)
	psw, _ := gopass.GetPasswdMasked()

	client.FileUpload(args[0], args[1], string(psw))
}

func NewFileDownloadCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "download <file id>",
		Short: "Download the any specific file you want",
		Long:  `Download command mean download file from the CESS networks based on fileId.`,

		Run: FileDownloadCommandFunc,
	}

	return cc
}

func FileDownloadCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)
	if len(args) == 0 {
		fmt.Printf("Please enter the fileid of the download file 'file download <fileid>'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}

	client.FileDownload(args[0])
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
	InitComponents(cmd)
	if len(args) == 0 {
		fmt.Printf("Please enter the fileid of the delete file'file delete <fileid>'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	client.FileDelete(args[0])

}

func NewFileDecryptCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "decrypt <file path>",
		Short: "Decrypt the any specific file again when you failed file decrypt first chance",
		Long:  `File decode means that if the file is not decrypted when you download it, it can be decrypt by this method.Please enter absolute path.`,

		Run: FileDecryptCommandFunc,
	}

	return cc
}

func FileDecryptCommandFunc(cmd *cobra.Command, args []string) {
	InitComponents(cmd)

	if len(args) == 0 {
		fmt.Printf("Please enter the path of the file to be decrypt'\n")
		os.Exit(conf.Exit_CmdLineParaErr)
	}
	client.FileDecrypt(args[0])
}
