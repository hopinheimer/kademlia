package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	bootnode := flag.Bool("bootnode", false, "initiate bootnode")
	port := flag.Int("port", 8081, "server port")

	flag.Parse()
	var srv *Server

	if *bootnode {
		fmt.Println("starting bootnode..... @port:", 8081)
		srv = NewServer(net.IPv4(127, 0, 0, 1), int32(8081))
	} else {
		srv = NewServer(net.IPv4(127, 0, 0, 1), int32(*port))
	}

	node := &Node{
		IP:   &srv.IP,
		Port: int(srv.Port),
	}
	rt := NewRoutingTable()
	node.rt = rt
	node.ID = node.toKadID(string(time.Now().UnixNano()))
	srv.LocalNode = node
	if err := srv.run(); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

	if !*bootnode {
		fmt.Println("connect with bootnode .... ")
		srv.ConnectWithBootnode()
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("graceful shutdown")
	srv.quit = true
	if srv.conn != nil {
		srv.conn.Close()
	}
}
