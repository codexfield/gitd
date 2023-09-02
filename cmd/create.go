/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"gitd/internal/storage"
	"gitd/internal/transport"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bnb-chain/greenfield-go-sdk/types"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	transport2 "github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/spf13/cobra"
)

const (
	DefaultBranchReferenceName = "refs/heads/main"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new repo on specify remote",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var url string
		if len(args) == 1 {
			url = args[0]
		} else {
			fmt.Println("Must specify a url, example: gnfd://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443/<reponame>")
			return
		}
		endpoint, err := transport2.NewEndpoint(url)
		if err != nil {
			fmt.Printf("New endpoint error: %s", err)
			return
		}
		//fmt.Println("Endpoint: ", endpoint.String())

		repoName, _ := strings.CutPrefix(endpoint.Path, "/")
		newStorage, err := storage.NewStorage(
			os.Getenv(transport.EnvChainID),
			"https://"+endpoint.Host+":"+strconv.Itoa(endpoint.Port),
			os.Getenv(transport.EnvPrivateKey),
			repoName,
		)
		if err != nil {
			fmt.Printf("New storage error: %s", err)
			return
		}
		_, err = newStorage.GnfdClient.HeadBucket(context.Background(), newStorage.GetBucketName())
		if err != nil {
			if strings.Contains(err.Error(), "No such bucket") {
				providers, err := newStorage.GnfdClient.ListStorageProviders(context.Background(), true)
				if err != nil {
					fmt.Println("list storage provider error: ", err)
					return
				}
				r := rand.New(rand.NewSource(time.Now().Unix()))

				if len(providers) > 0 {
					_, err := newStorage.GnfdClient.CreateBucket(context.Background(), newStorage.GetBucketName(), providers[r.Intn(len(providers))].OperatorAddress, types.CreateBucketOptions{})
					if err != nil {
						fmt.Println("create bucket error: ", err)
						return
					}
				}
			} else {
				fmt.Println("head bucket error: ", err)
				return
			}
		}

		_, err = git.InitWithOptions(newStorage, memfs.New(), git.InitOptions{DefaultBranch: DefaultBranchReferenceName})

		hash := plumbing.NewHashReference(DefaultBranchReferenceName, plumbing.Hash{})
		err = newStorage.SetReference(hash)
		if err != nil {
			fmt.Println("Set Reference ", DefaultBranchReferenceName, "error: ", err)
			return
		}
		fmt.Println("Created successfully!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
