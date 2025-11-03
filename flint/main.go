// /*
// Copyright © 2025 NAME HERE <EMAIL ADDRESS>

// */
package main

import (
	"time"

	"github.com/SuperTapood/Flint/cmd"
)

func main() {
	start := time.Now()
	cmd.Execute()
	println("done in " + (time.Since(start).String()))
}

// import (
// 	"log"
// 	"os"

// 	"github.com/SuperTapood/Flint/core"
// )

// func main() {
// 	data, err := os.ReadFile("bib.bin")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var stack = core.StackFromBinary(data)

// 	stack.Deploy()

// }
