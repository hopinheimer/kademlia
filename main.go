package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	node1 := NewKademlia(":8001")
	node2 := NewKademlia(":8002")

	var res bool
	client, err := rpc.Dial("tcp", node1.node.Addr)
	if err != nil {
		fmt.Println("Error dialing:", err)
		return
	}
	defer client.Close()
	err = client.Call("Kademlia.Ping", &node2.node, &res)
	if err != nil {
		fmt.Println("Error calling Ping:", err)
		return
	}
	fmt.Println("Ping result:", res)
}
