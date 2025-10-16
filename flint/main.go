package main

import (
	"log"
	"os"

	"github.com/SuperTapood/Flint/core/generated/common"
	"google.golang.org/protobuf/proto"
)

func main() {
	data, err := os.ReadFile("../bob.bin")
	if err != nil {
		log.Fatal(err)
	}

	var stack common.Stack_

	if err := proto.Unmarshal(data, &stack); err != nil {
		log.Fatal("failed to unmarshal:", err)
	}

	log.Printf("Loaded user: %+v\n", stack)

	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// resp, err := http.Get("https://192.168.49.2:8443/api/v1/namespaces/default/pods")

	// log.Print(resp)
	// log.Print(err)
}
