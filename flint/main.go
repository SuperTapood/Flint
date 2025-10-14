package main

import (
	"fmt"
	"log"
	"os"

	gen "github.com/SuperTapood/Flint/generated"
	"google.golang.org/protobuf/proto"
)

func main() {
	fmt.Println("hello world")

	data, err := os.ReadFile("../bob.bin")
	if err != nil {
		log.Fatal(err)
	}

	var person gen.Person

	if err := proto.Unmarshal(data, &person); err != nil {
		log.Fatal("failed to unmarshal:", err)
	}

	log.Printf("Loaded user: %+v\n", person)
}
