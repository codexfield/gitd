/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"golang.org/x/mod/module"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file contents to the index\n",
	Long:  `usage: git add [<options>] [--] <pathspec>...`,
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
		for _, arg := range args {
			if module.CheckFilePath(arg) == nil {
				_, err := w.Add(arg)
				if err != nil {
					fmt.Println("Add path failed, error: ", err)
					return
				}
			} else {
				fmt.Println("Check file path failed, error: ", err, "path", arg)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
