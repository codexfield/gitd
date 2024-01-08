/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"gitd/internal/storage"
	"gitd/internal/transport"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/bnb-chain/greenfield-go-sdk/types"
	accountmanager "github.com/codexfield/codex-contracts-go-sdk/account"
	"github.com/codexfield/codex-contracts-go-sdk/contracts/codexam"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	transport2 "github.com/go-git/go-git/v5/plumbing/transport"

	"github.com/spf13/cobra"
)

const (
	DefaultBranchReferenceName = "refs/heads/main"
	AccountManagerAddr         = "0xae5c57a7285602830aEA302f56e8Cf647a82F022"
	CodexBrand                 = "codex"
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
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			force = false
		}
		// get account id from cdoexfield account manager contract
		acc, err := accountmanager.NewAccount(os.Getenv(transport.EnvPrivateKey))
		if err != nil {
			fmt.Println("prepare account failed", "err", err)
			return
		}
		client, err := ethclient.Dial("https://data-seed-prebsc-1-s1.binance.org:8545/")
		if err != nil {
			fmt.Println("dial rpc failed", "err", err)
			return
		}
		codexAM, err := codexam.NewICodexAM(common.HexToAddress(AccountManagerAddr), client)
		if err != nil {
			fmt.Println("new codex account manager instance failed", "err", err)
			return
		}
		// Get Account ID
		accountID, err := codexAM.GetAccountId(&bind.CallOpts{}, acc.Address())
		if err != nil {
			fmt.Println("get account id failed", "err", err)
			return
		}
		if accountID.Cmp(big.NewInt(0)) <= 0 {
			fmt.Println("Unregister account. ")
			return
		}

		repoName, _ := strings.CutPrefix(endpoint.Path, "/")
		repoName = CodexBrand + "-" + accountID.String() + "-" + repoName
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

				if len(providers) > 0 {
					_, err := newStorage.GnfdClient.CreateBucket(context.Background(), newStorage.GetBucketName(), providers[0].OperatorAddress, types.CreateBucketOptions{})
					if err != nil {
						fmt.Println("create repo error: ", err)
						return
					}
				}
			} else {
				fmt.Println("create repo failed, error: ", err)
				return
			}
		} else {
			if !force {
				fmt.Println("Repo already exist. use --force if needs overwrite.")
				return
			}
		}

		_, err = git.InitWithOptions(newStorage, memfs.New(), git.InitOptions{DefaultBranch: DefaultBranchReferenceName})
		if err != nil {
			fmt.Println("Git init failed. Error: ", err)
		}

		hash := plumbing.NewHashReference(DefaultBranchReferenceName, plumbing.Hash{})
		err = newStorage.SetReference(hash)
		if err != nil {
			fmt.Println("Set Reference ", DefaultBranchReferenceName, "error: ", err)
			return
		}
		fmt.Println("Created successfully! Repo name: ", repoName)
	},
}

func init() {
	createCmd.Flags().BoolP("force", "f", false, "force create")
	// rootCmd.AddCommand(createCmd)
}
