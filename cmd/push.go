/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Update remote refs along with associated objects",
	Long:  `usage: git push [<options>] [<repository> [<refspec>...]]`,
	Run: func(cmd *cobra.Command, args []string) {
		r, err := git.PlainOpen("./")
		if err != nil {
			fmt.Println("Open repository failed, error: ", err)
			return
		}

		var (
			remoteName string
		)

		if len(args) >= 2 {
			remoteName = args[0]
		} else {
			fmt.Println("Please specify remote repository name and branch")
			return
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			force = false
		}

		err = r.Push(&git.PushOptions{
			RemoteName: remoteName,
			Force:      force,
		})
		if err != nil {
			fmt.Println("Repository push failed, error: ", err)
			return
		}
		fmt.Println("Repository push success.")
	},
}

func init() {
	pushCmd.Flags().BoolP("force", "f", false, "force updates")
	rootCmd.AddCommand(pushCmd)
}
