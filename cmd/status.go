/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the working tree status",
	Long:  `usage: git status`,
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
		head, err := r.Head()
		if err != nil || !head.Name().IsBranch() {
			fmt.Println("Get repo head error: ", err)
		}
		fmt.Println("On branch ", head.Name().Short())
		s, err := w.Status()
		if err != nil {
			fmt.Println("Get worktree status failed, error: ", err)
			return
		}
		if s.IsClean() {
			fmt.Println("nothing to commit, working tree clean")
		} else {
			fmt.Println("Changes not staged for commit:\n  (use \"git add <file>...\" to update what will be committed)\n  (use \"git restore <file>...\" to discard changes in working directory)\n")
			fmt.Println("  " + s.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
