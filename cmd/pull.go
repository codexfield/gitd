package cmd

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	cmd := &cobra.Command{
		Use:   "pull [remote]",
		Short: "Fetch from and integrate with another repository",
		Long: `Incorporate changes from a remote repository into the current branch.
If remote is not specified, the remote named origin is used.`,
		Args: cobra.MaximumNArgs(1),
		RunE: pullCmd,
	}
	rootCmd.AddCommand(cmd)
}

func pullCmd(cmd *cobra.Command, args []string) error {
	_, r, err := openRepo()
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	remote := "origin"
	if len(args) >= 1 {
		remote = args[0]
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: remote,
		Progress:   os.Stdout,
	})
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		fmt.Printf("%v\n", err)
		return nil
	}
	return err
}
