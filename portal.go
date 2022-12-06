package main

import (
	"cess-portal/command"

	"github.com/spf13/cobra"
)

const (
	Name        = "cessctl"
	Description = "Cess portal is used by entering the command line"
)

var (
	rootCmd = &cobra.Command{
		Use:        Name,
		Short:      Description,
		SuggestFor: []string{"cessctl"},
	}
	globalFlag = command.GlobalFlags{}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&globalFlag.ConfFilePath, "config", "c", "", "Custom configuration file path, requires absolute path")

	rootCmd.AddCommand(
		command.NewQueryCommand(),
		command.NewFileCommand(),
		command.NewSpaceCommand(),
		command.NewBucketCommand(),
	)
}
func Start() error {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	return rootCmd.Execute()
}

func main() {
	Start()
}
