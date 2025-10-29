/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [cluster-type]",
	Short: "list existing flint stacks",
	Long:  `list existing flint stacks`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list what dumbass")
		os.Exit(1)
	},
}

var listK8SCmd = &cobra.Command{
	Use:   "k8s [resource_type 1, ...]",
	Short: "list flint stacks in k8s",
	Long:  `list flint stacks in k8s`,
	Run:   listK8s,
	//Args:  cobra.ExactArgs(1),
}

var (
	token string
	api   string
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.AddCommand(listK8SCmd)

	listK8SCmd.Flags().SortFlags = false
	listK8SCmd.Flags().StringVarP(&token, "token", "t", "", "the token for the kubernetes cluster")
	// listK8SCmd.MarkFlagRequired("token")
	listK8SCmd.Flags().StringVarP(&api, "api", "a", "", "the api url of the kubernetes cluster")
	// listK8SCmd.MarkFlagRequired("api")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printList(conn base.Connection) {
	conn.List()
}

func listK8s(cmd *cobra.Command, args []string) {
	var bad = false

	if api == "" {
		bad = true
		fmt.Println("--api/-a was not specified")
	}

	if token == "" {
		bad = true
		fmt.Println("--token/-t was not specified")
	}

	if bad {
		fmt.Println("and a kubeconfig file was not found to be used instead")
		os.Exit(1)
	}

	conn := base.K8SConnection{
		Api:   api,
		Token: token,
	}

	printList(&conn)

}
