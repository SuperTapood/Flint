package main

import (
	"log"
	"os"

	"github.com/SuperTapood/Flint/core"
)

func main() {
	data, err := os.ReadFile("../bob.bin")
	if err != nil {
		log.Fatal(err)
	}

	var stack = core.StackFromBinary(data)
	log.Printf(stack.String())
	var dag, obj_map = stack.Synth()
	log.Printf("%v", &dag)
	log.Print("%v", obj_map)

	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// resp, err := http.Get("https://192.168.49.2:8443/api/v1/namespaces/default/pods")

	// log.Print(resp)
	// log.Print(err)
}
