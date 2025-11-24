package cmd

import (
	"fmt"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	deployMaxSecretNumber int
	deployForce           bool
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy a flint stack to the cloud",
	Long:  `deploy a flint stack to the cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		stack, conn, stackName := StackConnFromApp()
		objDag, objMap := stack.GetActual().Synth(stackName)
		added, removed, changed := conn.Diff(objMap, stack.GetActual().GetMetadata(), stackName)
		if !deployForce && len(added) == 0 && len(removed) == 0 && len(changed) == 0 {
			fmt.Println("empty changeset nothing to do")
			return
		}
		for _, name := range removed {
			unresource := base.Unresource{
				Name: name,
				ID:   uuid.New().String(),
			}
			objMap[unresource.GetID()] = &unresource
			objDag.AddVertexByID(unresource.GetID(), unresource.GetID())
		}
		conn.Deploy(objDag, objMap, stackName, stack.GetActual().GetMetadata(), deployMaxSecretNumber, true)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().SortFlags = false

	deployCmd.Flags().StringVarP(&app, "app", "a", "", "the app to synth the ")
	deployCmd.MarkFlagRequired("app")

	deployCmd.Flags().StringVarP(&dir, "dir", "d", ".", "the directory to run the app at")

	deployCmd.Flags().IntVar(&deployMaxSecretNumber, "history", 5, "the number of flint stacks you want remembered")
	deployCmd.Flags().BoolVarP(&deployForce, "force", "f", false, "if set, deploy even if there are no changes")
}
