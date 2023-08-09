package cmd

import (
	"fmt"
	"gitd/internal/transport"
	"github.com/spf13/cobra"
	"runtime"
)

const (
	Version = "0.0.1"
)

var rootCmd = &cobra.Command{
	Use:   "gitd",
	Short: "A gitd command tools based on Greenfield Decentralized Storage",
	Long: `usage: gitd [-v | --version] [-h | --help] [-C <path>] [-c <name>=<value>]
           [--exec-path[=<path>]] [--html-path] [--man-path] [--info-path]
           [-p | --paginate | -P | --no-pager] [--no-replace-objects] [--bare]
           [--git-dir=<path>] [--work-tree=<path>] [--namespace=<name>]
           [--super-prefix=<path>] [--config-env=<name>=<envvar>]
           <command> [<args>]`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		transport.InstallGreenfieldTransport()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			fmt.Println("gitd version", Version, "("+runtime.GOARCH+")")
		} else {
			fmt.Println(cmd.Help())
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode")
}
