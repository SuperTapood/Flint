package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SuperTapood/Flint/generated/k8s"
	"google.golang.org/protobuf/proto"
)

func main() {
	// create a Pod
	pod := &k8s.Pod{
		Name:  "nginx",
		Image: "nginx:latest",
		Ports: []int32{80, 443},
	}

	// // wrap it in K8STypes (oneof)
	k8stype := &pb.K8STypes{
		Type: &pb.K8STypes_Pod{Pod: pod},
	}

	// // create the stack
	// stack := &pb.XK8SStack{
	// 	Objects: []*pb.K8STypes{k8stype},
	// }

	// // marshal to binary
	// data, err := proto.Marshal(stack)
	// if err != nil {
	// 	panic(err)
	// }

	// // save it
	// if err := os.WriteFile("stack.bin", data, 0644); err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Wrote %d bytes to stack.bin\n", len(data))

	// // read back and verify
	// readData, _ := os.ReadFile("stack.bin")
	// var decoded pb.XK8SStack
	// if err := proto.Unmarshal(readData, &decoded); err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Decoded stack: %+v\n", decoded)
	data, err := os.ReadFile("../bob.bin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(data))

	log.Printf("%+v", data)

	var stack k8s.XK8SStack

	if err := proto.Unmarshal(data, &stack); err != nil {
		log.Printf("Loaded user: %+v\n", stack)
		log.Fatal("failed to unmarshal:", err)
	}

	log.Printf("Loaded user: %+v\n", stack)
}
