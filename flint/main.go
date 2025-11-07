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
