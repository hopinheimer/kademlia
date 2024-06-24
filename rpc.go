package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type Kademlia struct {
	node         Node
	routingTable *RoutingTable
}

func NewKademlia(addr string) *Kademlia {
	k := &Kademlia{
		node:         Node{ID: NewNodeID(), Addr: addr},
		routingTable: NewRoutingTable(),
	}
	rpc.Register(k)
	go k.listen()
	return k
}

func (k *Kademlia) listen() {
	listener, err := net.Listen("tcp", k.node.Addr)
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

func (k *Kademlia) Ping(req *Node, res *bool) error {
	k.routingTable.AddNode(*req)
	*res = true
	return nil
}

func (k *Kademlia) Store(req *StoreRequest, res *bool) error {

	*res = true
	return nil
}

func (k *Kademlia) FindNode(req *FindNodeRequest, res *FindNodeResponse) error {

	return nil
}

func (k *Kademlia) FindValue(req *FindValueRequest, res *FindValueResponse) error {
	return nil
}

type StoreRequest struct {
	Key   NodeID
	Value []byte
}

type FindNodeRequest struct {
	ID NodeID
}

type FindNodeResponse struct {
	Nodes []Node
}

type FindValueRequest struct {
	Key NodeID
}

type FindValueResponse struct {
	Value []byte
	Nodes []Node
}
