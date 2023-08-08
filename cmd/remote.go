/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"

	"github.com/spf13/cobra"
)

// remoteCmd represents the remote command
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("remote called")
	},
}

var addSubCmd = &cobra.Command{
	Use:   "add",
	Short: "usage: git remote add [<options>] <name> <url>\n",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := git.PlainOpen("./")
		if err != nil {
			fmt.Println("Open repository failed, error: ", err)
			return
		}
		if len(args) != 2 {
			fmt.Println("Parameter error: ", err)
			return
		}

		_, err = r.CreateRemote(&config.RemoteConfig{
			Name: args[0],
			URLs: []string{args[1]},
		})
		if err != nil {
			fmt.Println("Create remote error, error: ", err)
			return
		}
	},
}

var showSubCmd = &cobra.Command{
	Use:   "show",
	Short: "usage: git remote show [<options>] <name>\n",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := git.PlainOpen("./")
		if err != nil {
			fmt.Println("Open repository failed, error: ", err)
			return
		}
		remotes, err := r.Remotes()
		for _, remote := range remotes {
			fmt.Println(remote.String())
		}
	},
}

func init() {
	remoteCmd.AddCommand(addSubCmd)
	remoteCmd.AddCommand(showSubCmd)
	rootCmd.AddCommand(remoteCmd)
}
