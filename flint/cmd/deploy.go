/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	deployMaxSecretNumber int
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy a flint stack to the cloud",
	Long:  `deploy a flint stack to the cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		stack, conn, stack_name := StackConnFromApp()
		obj_dag, obj_map := stack.GetActual().Synth(stack_name)
		added, removed, changed := conn.GetActual().Diff(obj_map, stack_name)
		if len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
			fmt.Println("empty changeset nothing to do")
			return
		}
		conn.GetActual().Deploy(obj_dag, removed, obj_map, stack_name, stack.GetActual().GetMetadata(), deployMaxSecretNumber)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().SortFlags = false

	deployCmd.Flags().StringVarP(&app, "app", "a", "", "the app to synth the ")
	deployCmd.MarkFlagRequired("app")

	deployCmd.Flags().StringVarP(&dir, "dir", "d", ".", "the directory to run the app at")

	deployCmd.Flags().IntVar(&deployMaxSecretNumber, "history", 5, "the number of flint stacks you want remembered")
}
