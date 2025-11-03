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

func StackFromApp() *common.Stack {
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
		return &stack
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

	return common.StackFromBinary(data)
}
