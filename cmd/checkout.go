package cmd

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"

	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var checkCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Switch branches",
	Long:  `usage: git checkout [-b] `,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, r, err := openRepo()
		if err != nil {
			fmt.Println("Open repository failed, error: ", err)
			return err
		}
		w, err := r.Worktree()
		if err != nil {
			fmt.Println("Get worktree failed, error: ", err)
			return err
		}

		create, _ := cmd.Flags().GetBool("create")
		return w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(args[0]),
			Create: create,
			Force:  false,
			Keep:   true,
		})
	},
}

func init() {
	checkCmd.Flags().BoolP("create", "b", false, "Create branch before checkout")
	rootCmd.AddCommand(checkCmd)
}
