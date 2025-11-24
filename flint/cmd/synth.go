package cmd

import (
	"encoding/binary"
	"os"

	"github.com/spf13/cobra"
)

// synthCmd represents the synth command
var synthCmd = &cobra.Command{
	Use:   "synth",
	Short: "synth a stack into a stack file",
	Long:  `synth a stack into a stack file`,
	Run: func(cmd *cobra.Command, args []string) {
		stack, _, stack_name := StackConnFromApp()

		binary_synthed, _ := stack.GetActual().Synth(stack_name)

		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}

		defer file.Close()

		err = binary.Write(file, binary.LittleEndian, binary_synthed)
		if err != nil {
			panic(err)
		}
	},
}

var (
	filename string
)

func init() {
	rootCmd.AddCommand(synthCmd)

	synthCmd.Flags().StringVarP(&filename, "file", "f", "stack.out", "specify the file name to write the synth output to")
	synthCmd.Flags().StringVarP(&app, "app", "a", "", "the app to synth the ")
	synthCmd.MarkFlagRequired("app")
	synthCmd.Flags().StringVarP(&dir, "dir", "d", ".", "the directory to run the app at")
}
