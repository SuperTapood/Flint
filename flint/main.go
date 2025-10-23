package main

import (
	"log"
	"os"

	"github.com/SuperTapood/Flint/core"
)

func main() {
	data, err := os.ReadFile("bib.bin")
	if err != nil {
		log.Fatal(err)
	}

	var stack = core.StackFromBinary(data)
	log.Print(stack.String())
	var dag, obj_map = stack.Synth()
	log.Printf("%v", dag)
	log.Printf("%v", obj_map)

	stack.Deploy()

}
