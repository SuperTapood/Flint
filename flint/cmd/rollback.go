package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rollbackTargetRevision int
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "rollback your stack to an earlier",
	Long:  `rollback your stack to an earlier`,
	Run: func(cmd *cobra.Command, args []string) {
		stack, conn, stack_name := StackConnFromApp()
		revision := rollbackTargetRevision
		if rollbackTargetRevision < 0 {
			revision = conn.GetActual().GetCurrentRevision(stack_name)
			revision += rollbackTargetRevision
		}
		conn.GetActual().Rollback(stack_name, revision, stack.GetActual().GetMetadata())
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)

	rollbackCmd.Flags().StringVarP(&app, "app", "a", "", "the app to destroy ")
	rollbackCmd.MarkFlagRequired("app")
	rollbackCmd.Flags().StringVarP(&dir, "dir", "d", ".", "the directory to run the app at")
	rollbackCmd.Flags().IntVarP(&rollbackTargetRevision, "number", "n", 0, "the revision of the stack to rollback to")
	rollbackCmd.MarkFlagRequired("number")
}
