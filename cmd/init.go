/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create an empty Git repository or reinitialize an existing one",
	Long: `usage: gitd init [-q | --quiet] [--bare] [--template=<template-directory>]
                [--separate-git-dir <git-dir>] [--object-format=<format>]
                [-b <branch-name> | --initial-branch=<branch-name>]
                [--shared[=<permissions>]] [<directory>]`,
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) == 0 {
			path = "./"
		} else {
			path = args[0]
		}

		fs := osfs.New(path)
		dot, _ := fs.Chroot(".git")
		_, err := git.Init(filesystem.NewStorage(dot, cache.NewObjectLRUDefault()), nil)
		if err != nil {
			fmt.Println("init repository error: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
