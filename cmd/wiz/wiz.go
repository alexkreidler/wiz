package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tim15/wiz/cli/flags"
	"os"
)

func newWizCmd() *cobra.Command {
	opts := flags.NewClientOptions()
	var cmd = &cobra.Command{
		Use:           "wiz [OPTIONS] COMMAND",
		Short:         "An ML platform and package manager",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Version {
				showVersion()
			}
			return nil
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&opts.Version, "version", "v", false, "Print version information and quit")
	// flags
	return cmd
}

func showVersion() {
	fmt.Printf("Wiz CLI version %s\n", "0.0.1-dev")
}

func main() {
	RootCmd := newWizCmd()
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
