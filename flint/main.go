package main

import (
	"log"
	"os"

	"github.com/SuperTapood/Flint/generated/k8s"
	"google.golang.org/protobuf/proto"
)

func main() {
	data, err := os.ReadFile("../bob.bin")
	if err != nil {
		log.Fatal(err)
	}

	var stack k8s.XK8SStack

	if err := proto.Unmarshal(data, &stack); err != nil {
		log.Printf("Loaded user: %+v\n", stack)
		log.Fatal("failed to unmarshal:", err)
	}

	log.Printf("Loaded user: %+v\n", stack)

	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// resp, err := http.Get("https://192.168.49.2:8443/api/v1/namespaces/default/pods")

	// log.Print(resp)
	// log.Print(err)
}
