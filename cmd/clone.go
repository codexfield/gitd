/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/bnb-chain/greenfield-go-sdk/types"
	"github.com/codexfield/gitd/storage"
	"github.com/codexfield/gitd/transport"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	transport2 "github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/spf13/cobra"
)

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone a repository into a new directory\n",
	Long:  `usage: git clone <repo> [<dir>]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("no URL provided")
		}

		url := args[0]
		// parse the endpoint
		endpoint, err := transport2.NewEndpoint(url)
		if err != nil {
			return err
		}

		repoName, _ := strings.CutPrefix(endpoint.Path, "/")
		backend, err := storage.NewStorage(
			os.Getenv(transport.EnvChainID),
			"https://"+endpoint.Host+":"+strconv.Itoa(endpoint.Port),
			os.Getenv(transport.EnvPrivateKey),
			repoName,
		)
		if err != nil {
			return err
		}

		dir := ""

		if len(args) >= 2 {
			dir = args[1]
		} else {
			dir = cloneDir(endpoint.Path)
		}

		fmt.Printf("Clone into %q...\n", dir)

		result, err := backend.GnfdClient.ListObjects(context.Background(), repoName, types.ListObjectsOptions{})
		if err != nil {
			return err
		}
		if len(result.Objects) == 1 && result.Objects[0].ObjectInfo.ObjectName == "refs/HEAD" {
			fmt.Println("warning: You appear to have cloned an empty repository.")

			ref, err := backend.Reference(plumbing.HEAD)
			if err != nil {
				return err
			}
			fs := osfs.New(dir)
			dot, _ := fs.Chroot(".git")
			r, err := git.InitWithOptions(filesystem.NewStorage(dot, cache.NewObjectLRUDefault()), memfs.New(), git.InitOptions{
				DefaultBranch: ref.Target(),
			})
			if err != nil {
				return err
			}
			_, err = r.CreateRemote(&config.RemoteConfig{
				Name: "origin",
				URLs: []string{url},
			})
			return err
		}

		_, err = git.PlainClone(dir, false, &git.CloneOptions{
			URL:      endpoint.String(),
			Progress: os.Stdout,
		})
		return err
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

// cloneDir return the directory where newly cloned repository will be stored
// based on the full path in git URL. For example, it converts
// fhs/gig to gig, and fhs/gig.git to gig.
func cloneDir(dir string) string {
	d := path.Base(dir)
	if suf := ".git"; strings.HasSuffix(d, suf) && len(d) > len(suf) {
		d = d[:len(d)-len(suf)]
	}
	return d
}
