/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/SuperTapood/Flint/core"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy a flint stack to the cloud",
	Long:  `deploy a flint stack to the cloud`,
	Run:   deploy,
}

var (
	app string
	dir string
)

func init() {
	rootCmd.AddCommand(deployCmd)

	deployCmd.Flags().StringVarP(&app, "app", "a", "", "the app to synth the ")
	deployCmd.MarkFlagRequired("app")

	deployCmd.Flags().StringVarP(&dir, "dir", "d", ".", "the directory to run the app at")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deploy(cmd *cobra.Command, args []string) {
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("Pipe error:", err)
		return
	}

	parts := strings.Fields(app + " 3")
	command := exec.Command(parts[0], parts[1:]...)

	command.Dir = dir
	command.ExtraFiles = []*os.File{w}
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	if err := command.Start(); err != nil {
		fmt.Println(err)
	}
	w.Close()

	var data []byte

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		data, err = io.ReadAll(r)
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		//fmt.Println("Received binary:", string(data)) // Process bytes as needed

	}()

	command.Wait()
	wg.Wait()
	r.Close()

	stack := core.StackFromBinary(data)

	stack.Deploy()

	// fmt.Println(string(out))

	// var obj proto.Message

	// err = proto.Unmarshal(out, obj)
	// fmt.Println(err)
	// fmt.Println(obj)
}
