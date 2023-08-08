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
		err = r.Push(&git.PushOptions{
			RemoteName: "origin",
		})
		if err != nil {
			fmt.Println("Repository push failed, error: ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

}
