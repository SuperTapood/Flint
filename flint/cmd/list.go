package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/SuperTapood/Flint/core/generated/k8s"
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
	k8s_token string
	k8s_api   string
)

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.AddCommand(listK8SCmd)

	listK8SCmd.Flags().SortFlags = false
	listK8SCmd.Flags().StringVarP(&k8s_token, "token", "t", "", "the token for the kubernetes cluster")
	listK8SCmd.Flags().StringVarP(&k8s_api, "api", "a", "", "the api url of the kubernetes cluster")
}

func listK8s(cmd *cobra.Command, args []string) {
	var bad = false

	if k8s_api == "" {
		bad = true
		fmt.Println("--api/-a was not specified")
	}

	if k8s_token == "" {
		bad = true
		fmt.Println("--token/-t was not specified")
	}

	if bad {
		fmt.Println("and a kubeconfig file was not found to be used instead")
		os.Exit(1)
	}

	conn := k8s.K8S_Connection{
		Api:   k8s_api,
		Token: k8s_token,
	}

	deployments := conn.List()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "Name\tAge\tStatus\tRevision")

	for i := range deployments {
		deployment := &deployments[i]
		fmt.Fprintln(w, deployment.Name+"\t"+deployment.Age+"\t"+deployment.Status+"\t"+strconv.Itoa(int(deployment.Revision)))
	}
	w.Flush()

}
