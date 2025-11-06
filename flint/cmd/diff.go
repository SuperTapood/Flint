/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "display the difference between the given stack and the existing one (if exists)",
	Long:  `display the difference between the given stack and the existing one (if exists)`,
	Run: func(cmd *cobra.Command, args []string) {
		// stack, conn, stack_name := StackConnFromApp()
		// _, obj_map := stack.GetActual().Synth(stack_name)
		// conn.Diff(obj_map, stack_name)
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().SortFlags = false

	diffCmd.Flags().StringVarP(&app, "app", "a", "", "the app to synth the ")
	diffCmd.MarkFlagRequired("app")
	diffCmd.Flags().StringVarP(&dir, "dir", "d", ".", "the directory to run the app at")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// diffCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// diffCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
