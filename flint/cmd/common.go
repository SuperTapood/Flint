package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

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
	r, w, err := os.Pipe()
	if err != nil {
		panic("Pipe error:" + err.Error())
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

	// return common.StackFromBinary(data)

	var stack common.Stack
	err = proto.Unmarshal(data, &stack)
	if err != nil {
		panic(err)
	}

	return stack.GetStack(), stack.GetConnection(), stack.GetName()
}
