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
	Short: "usage: git remote [-v | --verbose]\n",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		showSubCmd.Run(cmd, args)
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
		verbose, _ := cmd.Flags().GetBool("verbose")
		remotes, err := r.Remotes()
		for _, remote := range remotes {
			if verbose {
				fmt.Println(remote.String())
			} else {
				fmt.Println(remote.Config().Name)
			}
		}
	},
}

var removeSubCmd = &cobra.Command{
	Use:   "remove",
	Short: "usage: git remote remove <name>\n",
	Run: func(cmd *cobra.Command, args []string) {
		r, err := git.PlainOpen("./")
		if err != nil {
			fmt.Println("Open repository failed, error: ", err)
			return
		}
		if len(args) != 1 {
			fmt.Println("Too much args, error: ", err)
			return
		}
		err = r.DeleteRemote(args[0])
		if err != nil {
			fmt.Println("Delete remote failed, error: ", err)
			return
		}
	},
}

func init() {
	remoteCmd.AddCommand(addSubCmd)
	remoteCmd.AddCommand(showSubCmd)
	remoteCmd.AddCommand(removeSubCmd)
	rootCmd.AddCommand(remoteCmd)

	commitCmd.Flags().BoolP("verbose", "v", false, "amend previous commit")
}
