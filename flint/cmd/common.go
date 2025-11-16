package cmd

import (
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SuperTapood/Flint/core/generated/common"
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
func StackConnFromApp() (*common.StackTypes, *common.ConnectionTypes, string) {
	if _, err := os.Stat(app); err == nil {
		data, err := os.ReadFile(app)
		if err != nil {
			panic(err)
		}
		var stack common.Stack
		err = proto.Unmarshal(data, &stack)
		if err != nil {
			panic(err)
		}
		return stack.GetStack(), stack.GetConnection(), stack.GetName()
	}

	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "my-socket.sock")

	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	defer os.Remove(socketPath)

	parts := strings.Fields(app + " " + socketPath)
	command := exec.Command(parts[0], parts[1:]...)

	command.Dir = dir
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout

	if err := command.Start(); err != nil {
		panic(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024*1024)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}

	data := buf[:n]

	var stack common.Stack
	err = proto.Unmarshal(data, &stack)
	if err != nil {
		panic(err)
	}

	return stack.GetStack(), stack.GetConnection(), stack.GetName()
}
