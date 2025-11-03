/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/binary"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

// synthCmd represents the synth command
var synthCmd = &cobra.Command{
	Use:   "synth",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		stack := StackFromApp()

		binary_synthed, err := proto.Marshal(stack)

		if err != nil {
			panic(err)
		}

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// synthCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// synthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
