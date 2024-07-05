package main

import (
	"crypto/sha1"
	"net"
)

type Node struct {
	ID [20]byte

	IP   *net.IP
	Port int
	rt   *RoutingTable
}

func (n *Node) toKadID(key string) [20]byte {
	bs := sha1.Sum([]byte(key))
	return bs
}
