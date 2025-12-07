package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SuperTapood/Flint/core/generated/general"
	"google.golang.org/protobuf/proto"
)

// common useful variables
var (
	app     string
	dir     string
	noColor bool
)

// ANSI color codes
const (
	colorReset    = "\x1b[0m"  // Reset all attributes
	colorRed      = "\x1b[31m" // Red text
	colorGreen    = "\x1b[32m" // Green text
	colorYellow   = "\x1b[33m"
	unchagedColor = "\033[38;5;181m"
)

/*
Load stack from either a runnable program or a synthed file.

Returns:
  - StackTypes* - The abstract protobuf stack (needs to be `GetActual()`-ed to get a the actual useable `StackType`)
  - ConnectionTypes* - The abstract protobuf connection (needs to be `GetActual()`-ed to get a the actual useable `ConnectionType`)
  - string - The name of the stack. This value is inaccessible later on.
*/
func StackConnFromApp() (*general.StackTypes, *general.ConnectionTypes, string) {
	if _, err := os.Stat(app); err == nil {
		data, err := os.ReadFile(app)
		if err != nil {
			fmt.Println("could not read synthed app file")
			fmt.Println(err)
			os.Exit(-1)
		}
		var stack general.Stack
		err = proto.Unmarshal(data, &stack)
		if err != nil {
			fmt.Println("could not unmarshal stack from synthed app file")
			fmt.Println(err)
			os.Exit(-1)
		}
		return stack.GetStack(), stack.GetConnection(), stack.GetName()
	}

	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "my-socket.sock")

	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("could not open socket")
		fmt.Println(err)
		os.Exit(-1)
	}
	defer listener.Close()
	defer os.Remove(socketPath)

	parts := strings.Fields(app + " " + socketPath)
	command := exec.Command(parts[0], parts[1:]...)

	command.Dir = dir
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	if err := command.Run(); err != nil {
		fmt.Println("running command failed")
		fmt.Println(err)
		os.Exit(-1)
	}

	if command.ProcessState.ExitCode() != 0 {
		fmt.Printf("running '%v' failed with exit code %v\n\n", app, command.ProcessState.ExitCode())
		os.Exit(1)
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("failed to get data from the app")
		fmt.Println(err)
		os.Exit(-1)
	}
	defer conn.Close()

	buf := make([]byte, 1024*1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("failed to read app data")
		fmt.Println(err)
		os.Exit(-1)
	}

	data := buf[:n]

	var stack general.Stack
	err = proto.Unmarshal(data, &stack)
	if err != nil {
		fmt.Println("failed to unmarshal stack from app")
		fmt.Println(err)
		os.Exit(-1)
	}

	return stack.GetStack(), stack.GetConnection(), stack.GetName()
}
