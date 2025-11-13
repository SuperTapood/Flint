package cmd

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/SuperTapood/Flint/core/generated/common"
	"google.golang.org/protobuf/proto"
)

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
	// r, w, err := os.Pipe()
	// if err != nil {
	// 	panic("Pipe error:" + err.Error())
	// }

	// lis, err := net.Listen("tcp", "127.0.0.1:51001")
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	// // 2. Get the listener's actual address (e.g., "localhost:54321")
	// serverAddr := lis.Addr().String()
	// log.Printf("Go server listening on: %s", serverAddr)

	// // 3. Create the gRPC server
	// s := grpc.NewServer()

	// // 4. Start the Python child process
	// log.Println("Starting Python client subprocess...")

	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "my-socket.sock")

	// Remove any existing socket
	os.Remove(socketPath)

	// Create Unix socket listener
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
		fmt.Println(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024*1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Read error:", err)
	}

	data := buf[:n]

	var stack common.Stack
	err = proto.Unmarshal(data, &stack)
	if err != nil {
		panic(err)
	}

	return stack.GetStack(), stack.GetConnection(), stack.GetName()
}
