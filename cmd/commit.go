/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Long:  `usage: git commit -m <msg>`,
	Run: func(cmd *cobra.Command, args []string) {
		r, err := git.PlainOpen("./")
		if err != nil {
			fmt.Println("Open repository failed, error: ", err)
			return
		}
		w, err := r.Worktree()
		if err != nil {
			fmt.Println("Get worktree failed, error: ", err)
			return
		}
		msg, _ := cmd.Flags().GetString("message")
		_, err = w.Commit(msg, &git.CommitOptions{})
		if err != nil {
			return
		}
	},
}

func init() {
	commitCmd.Flags().StringP("message", "m", "", "Commit message")
	rootCmd.AddCommand(commitCmd)
}
