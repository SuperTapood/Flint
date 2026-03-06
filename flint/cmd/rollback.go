/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strconv"

	"github.com/SuperTapood/Flint/core/generated/general"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

var (
	revision int32
)

// rollbackCmd represents the rollback command
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "rollback a stack to a previous state",
	Long:  `rollback a stack to a previous state`,
	Run: func(cmd *cobra.Command, args []string) {
		_, conn, stackName := StackConnFromApp()
		_, _, _, _, version := conn.GetActual().GetLatestRevision(stackName)
		if revision < 1 {
			revision = version + int32(revision)
		}
		revisions := conn.GetActual().GetRevisions()
		rev := revisions[stackName+"-"+strconv.FormatInt(int64(revision), 10)]
		decodedStack := rev["stack"].([]byte)
		var oldStack general.StackTypes
		err := proto.Unmarshal(decodedStack, &oldStack)
		if err != nil {
			panic(err)
		}
		Deploy(&oldStack, conn, stackName)
	},
}

func init() {
	rootCmd.AddCommand(rollbackCmd)

	rollbackCmd.Flags().StringVarP(&app, "app", "a", "", "the app to destroy ")
	rollbackCmd.MarkFlagRequired("app")
	rollbackCmd.Flags().StringVarP(&dir, "dir", "d", ".", "the directory to run the app at")
	rollbackCmd.Flags().Int32VarP(&revision, "revision", "r", 0, "the revision you want to rollback to (if a negative number is suppied, it represents how many revisions back you want to go)")
}
