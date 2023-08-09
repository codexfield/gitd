/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	transport2 "github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/spf13/cobra"
	"strings"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a repository into a new directory\n",
	Long:  `usage: git clone <repo> [<dir>]`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			path string
			repo string
		)

		if len(args) == 0 {
			fmt.Println(cmd.Help())
		} else if len(args) == 1 {
			repo = args[0]
			path = "./"
		} else {
			repo = args[0]
			path = args[1]
		}

		// parse the endpoint
		endpoint, err := transport2.NewEndpoint(repo)
		if err != nil {
			fmt.Printf("The Repo URL error: %s", err)
			return
		}

		repoName, found := strings.CutPrefix(endpoint.Path, "/")
		if !found {
			fmt.Println("The repo url error: the repo name not exist")
			return
		}
		fs := osfs.New(path + repoName)
		dot, _ := fs.Chroot(".git")
		_, err = git.Clone(filesystem.NewStorage(dot, cache.NewObjectLRUDefault()), nil, &git.CloneOptions{
			URL: endpoint.String(),
		})
		if err != nil {
			fmt.Println("Clone repo failed, error: ", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}
