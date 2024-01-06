/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file contents to the index\n",
	Long:  `usage: git add [<options>] [--] <pathspec>...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			fmt.Printf("Nothing specified, nothing added.")
		}

		root, r, err := openRepo()
		if err != nil {
			return err
		}

		w, err := r.Worktree()
		if err != nil {
			return err
		}

		for _, arg := range args {
			a, err := repoRelPath(root, arg)
			if err != nil {
				return err
			}

			_, err = w.Add(a)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
