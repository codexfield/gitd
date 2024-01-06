/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show commit logs",
	Long:  `usage: git log`,
	Run: func(cmd *cobra.Command, args []string) {
		_, r, err := openRepo()
		if err != nil {
			fmt.Println("Open repository failed, error: ", err)
			return
		}
		commitIter, err := r.Log(&git.LogOptions{})
		if err != nil {
			if err == plumbing.ErrReferenceNotFound {
				fmt.Println("Empty repository!")
			} else {
				fmt.Println("Log failed, error: ", err)
			}
			return
		}
		i := 0
		pageSize := 10
		for {
			commit, err := commitIter.Next()
			if err != nil {
				return
			}
			fmt.Println(commit.String())
			if (i+1)%pageSize == 0 {
				fmt.Println("Press Enter to continue or 'q' to quit...")
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.TrimSpace(input)

				if input == "q" {
					break
				}
			}
			i++
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
