package cmd

import (
	"github.com/spf13/cobra"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy a stack",
	Long:  `destroy a stack`,
	Run: func(cmd *cobra.Command, args []string) {
		stack, conn, stackName := StackConnFromApp()
		conn.Destroy(stackName, stack.GetActual().GetMetadata())
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)

	destroyCmd.Flags().StringVarP(&app, "app", "a", "", "the app to destroy ")
	destroyCmd.MarkFlagRequired("app")
	destroyCmd.Flags().StringVarP(&dir, "dir", "d", ".", "the directory to run the app at")
}
