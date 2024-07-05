package main

import (
	"container/list"
	"fmt"
	"sync"
)

const (
	IDLength = 20
	K        = 20
)

type Bucket struct {
	list *list.List
}

type RoutingTable struct {
	buckets []*Bucket
	mutex   sync.Mutex
}

func NewBucket() *Bucket {
	return &Bucket{
		list: list.New(),
	}
}

func NewRoutingTable() *RoutingTable {
	rt := &RoutingTable{
		buckets: make([]*Bucket, IDLength*8),
	}
	for i := range rt.buckets {
		rt.buckets[i] = NewBucket()
	}
	return rt
}

func (n *Node) AddNode(node *Node) {
	n.rt.mutex.Lock()
	defer n.rt.mutex.Unlock()

	bucketIndex := PrefixLen(XOR(n.ID, node.ID))
	bucket := n.rt.buckets[bucketIndex]

	for e := bucket.list.Front(); e != nil; e = e.Next() {
		if e.Value.(*Node).ID == node.ID {

			bucket.list.MoveToFront(e)
			return
		}
	}

	if bucket.list.Len() < K {
		bucket.list.PushFront(node)
	} else {
		//TODO:
		fmt.Println("eviction not implemented")
	}
}

func (n *Node) FindClosest(target [20]byte, count int) []*Node {
	n.rt.mutex.Lock()
	defer n.rt.mutex.Unlock()

	result := []*Node{}
	bucketIndex := PrefixLen(XOR(n.ID, target))
	bucket := n.rt.buckets[bucketIndex]

	for e := bucket.list.Front(); e != nil && len(result) < count; e = e.Next() {
		result = append(result, e.Value.(*Node))
	}

	for i := 1; len(result) < count && (bucketIndex-i >= 0 || bucketIndex+i < len(n.rt.buckets)); i++ {
		if bucketIndex-i >= 0 {
			for e := n.rt.buckets[bucketIndex-i].list.Front(); e != nil && len(result) < count; e = e.Next() {
				result = append(result, e.Value.(*Node))
			}
		}
		if bucketIndex+i < len(n.rt.buckets) {
			for e := n.rt.buckets[bucketIndex+i].list.Front(); e != nil && len(result) < count; e = e.Next() {
				result = append(result, e.Value.(*Node))
			}
		}
	}

	return result
}
