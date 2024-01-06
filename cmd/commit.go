/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Record changes to the repository",
	Long:  `usage: git commit -m <msg>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, r, err := openRepo()
		if err != nil {
			fmt.Println("Open repository failed, error: ", err)
			return err
		}
		w, err := r.Worktree()
		if err != nil {
			fmt.Println("Get worktree failed, error: ", err)
			return err
		}
		status, err := w.Status()
		if err != nil {
			return err
		}
		all, _ := cmd.Flags().GetBool("all")
		if nothingInStaging(status) {
			if !all {
				return fmt.Errorf(`nothing to commit (use "gig add")`)
			}
			if nothingInWorktree(status) {
				return fmt.Errorf(`nothing to commit, working tree clean`)
			}
		}
		commitMsg, _ := cmd.Flags().GetString("message")

		if commitMsg == "" {
			mfile := filepath.Join(root, ".git", "COMMIT_EDITMSG")
			err := os.WriteFile(mfile, []byte(emptyCommit), 0644)
			if err != nil {
				return err
			}
			err = editFile(mfile)
			if err != nil {
				return fmt.Errorf("editor failed: %v", err)
			}
			msg, err := readCommit(mfile)
			if err != nil {
				return err
			}
			if msg == "" {
				return fmt.Errorf("aborting commit due to empty commit message")
			}
			commitMsg = msg
		}
		_, err = w.Commit(commitMsg, &git.CommitOptions{
			All: all,
		})
		return err
	},
}

func init() {
	commitCmd.Flags().StringP("message", "m", "", "Commit message")
	commitCmd.Flags().BoolP("all", "a", false, "Stage modified/deleted files before commit")
	rootCmd.AddCommand(commitCmd)
}

func nothingInStaging(s git.Status) bool {
	for _, status := range s {
		switch status.Staging {
		case git.Unmodified, git.Untracked:
		default:
			return false
		}
	}
	return true
}

func nothingInWorktree(s git.Status) bool {
	for _, status := range s {
		switch status.Worktree {
		case git.Unmodified, git.Untracked:
		default:
			return false
		}
	}
	return true
}

func editFile(filename string) error {
	args := strings.Fields(preferredEditor())
	if len(args) == 0 {
		panic("internal error: empty editor")
	}
	args = append(args, filename)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func readCommit(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var msg []byte
	for scanner.Scan() {
		b := append(scanner.Bytes(), '\n')
		if b[0] != '#' {
			msg = append(msg, b...)
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	msg = bytes.TrimFunc(msg, func(r rune) bool { return r == '\n' })
	if len(msg) > 0 {
		msg = append(msg, '\n')
	}
	return string(msg), nil
}

func preferredEditor() string {
	for _, name := range []string{
		"GIT_EDITOR",
		"VISUAL",
		"EDITOR",
	} {
		if e := os.Getenv(name); e != "" {
			return e
		}
	}
	return defaultEditor()
}

var unknownUserMsg = `Author identity unknown

*** Please tell me who you are.

To set your account's default identity, write a configuration file like
the following to .git/config file or the global <HOME>/.gitconfig file:

	[user]
	email = gitster@example.com
	name = Junio C Hamano

Alternatively, set the GIT_AUTHOR_NAME and GIT_AUTHOR_EMAIL environment
variables instead.
`

var emptyCommit = `
# Please enter the commit message for your changes. Lines starting
# with '#' will be ignored, and an empty message aborts the commit.
`

func defaultEditor() string {
	for _, e := range []string{
		"vim",
		"vi",
		"nano",
	} {
		if p, err := exec.LookPath(e); err == nil {
			return p
		}
	}
	return "ed"
}
