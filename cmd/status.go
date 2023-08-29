/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

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
		if err != nil {
			if errors.Is(err, plumbing.ErrReferenceNotFound) {
				head, err = r.Storer.Reference(plumbing.HEAD)
				if err != nil {
					fmt.Println("Get empty repo head error: ", err)
					return
				}
				fmt.Println("On branch ", head.Target().Short())
			} else {
				fmt.Println("Get repo head error: ", err)
				return
			}
		} else {
			fmt.Println("On branch ", head.Name().Short())
		}

		s, err := w.Status()
		if err != nil {
			fmt.Println("Get worktree status failed, error: ", err)
			return
		}
		if s.IsClean() {
			fmt.Println("nothing to commit, working tree clean")
		} else {
			fmt.Println("Changes not staged for commit:\n  (use \"git add <file>...\" to update what will be committed)")
			fmt.Println("  " + s.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
