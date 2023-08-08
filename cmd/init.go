/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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
		fmt.Println(cmd.Help())
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
