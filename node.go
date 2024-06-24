package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	IDLength   = 160 // Length of node ID in bits
	BucketSize = 20  // Number of nodes per bucket
)

type NodeID [IDLength / 8]byte

type Node struct {
	ID   NodeID
	Addr string
}

type Bucket struct {
	nodes []Node
}

type RoutingTable struct {
	buckets [IDLength]*Bucket
}

func NewNodeID() NodeID {
	var id NodeID
	rand.Read(id[:])
	return id
}

func (n NodeID) XOR(other NodeID) NodeID {
	var result NodeID
	for i := 0; i < len(n); i++ {
		result[i] = n[i] ^ other[i]
	}
	return result
}

func (n NodeID) Distance(other NodeID) *big.Int {
	xorResult := n.XOR(other)
	return new(big.Int).SetBytes(xorResult[:])
}

func NewRoutingTable() *RoutingTable {
	rt := &RoutingTable{}
	for i := range rt.buckets {
		rt.buckets[i] = &Bucket{}
	}
	return rt
}

func (rt *RoutingTable) AddNode(node Node) {
	distance := rt.buckets[0].nodes[0].ID.Distance(node.ID)
	bucketIndex := distance.BitLen() - 1
	rt.buckets[bucketIndex].AddNode(node)
}

func (b *Bucket) AddNode(node Node) {
	if len(b.nodes) < BucketSize {
		b.nodes = append(b.nodes, node)
	} else {
		fmt.Print("some random case")
		// Handle bucket full scenario (e.g., eviction or pinging existing nodes)
	}
}
