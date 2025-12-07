package cmd

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// synthCmd represents the synth command
var synthCmd = &cobra.Command{
	Use:   "synth",
	Short: "synth a stack into a stack file",
	Long:  `synth a stack into a stack file`,
	Run: func(cmd *cobra.Command, args []string) {
		stack, _, stackName := StackConnFromApp()

		binarySynthed, _ := stack.GetActual().Synth(stackName)

		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("couldn't create the file to synth to")
			fmt.Println(err)
			os.Exit(1)
		}

		defer file.Close()

		err = binary.Write(file, binary.LittleEndian, binarySynthed)
		if err != nil {
			fmt.Println("couldn't write to the file when synthing")
			fmt.Println(err)
			os.Exit(1)
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
