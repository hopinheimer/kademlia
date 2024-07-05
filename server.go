package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

type UDPMessage struct {
	ID            [20]byte
	LocalAddress  *net.IP
	RemoteAddress *net.IP
	LocalPort     int32
	RemotePort    int32
	MessageClass  MessageType
	data          string
}

type MessageType string

const (
	Ping        MessageType = "ping"
	Pong        MessageType = "pong"
	Add         MessageType = "add"
	Drop        MessageType = "drop"
	SendMessage MessageType = "send"
	RecvMessage MessageType = "recv"
)

type Server struct {
	IP       net.IP
	Port     int32
	conn     *net.UDPConn
	quit     bool
	recvChan chan *UDPMessage
	sendChan chan *UDPMessage

	LocalNode *Node
}

func NewServer(ip net.IP, port int32) *Server {
	return &Server{
		IP:       ip,
		Port:     port,
		sendChan: make(chan *UDPMessage, 10),
		recvChan: make(chan *UDPMessage, 10),
	}
}

func (srv *Server) run() error {
	src := &net.UDPAddr{IP: srv.IP, Port: int(srv.Port), Zone: ""}
	conn, err := net.ListenUDP("udp", src)
	if err != nil {
		return err
	}

	srv.quit = false
	srv.conn = conn
	srv.sendChan = make(chan *UDPMessage, 10)

	go srv.readLoop()
	go srv.sendLoop()
	go srv.processMessage()

	return nil
}

func (srv *Server) Serialize(msg *UDPMessage) []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(msg)
	if err != nil {
		fmt.Println("Error encoding UDPMessage:", err)
		return nil
	}

	return buf.Bytes()
}

func (srv *Server) Deserialize(data []byte) (*UDPMessage, error) {
	var msg UDPMessage
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)

	err := decoder.Decode(&msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (srv *Server) sendLoop() {
	func() {
		for sendMsg := range srv.sendChan {
			dst := &net.UDPAddr{IP: *sendMsg.RemoteAddress, Port: int(sendMsg.RemotePort)}
			fmt.Println("sending connect request to bootnode")
			_, err := srv.conn.WriteToUDP(srv.Serialize(sendMsg), dst)
			if err != nil {
				fmt.Print("error")
			}
		}

	}()
}

func (srv *Server) readLoop() {

	go func() {
		var buf [1500]byte

		for {
			n, _, err := srv.conn.ReadFromUDP(buf[0:])
			if err != nil {
				if srv.quit {
					break
				}
				continue
			}

			data := make([]byte, n)
			copy(data, buf[:n])
			msg, err := srv.Deserialize(data)
			if err != nil {
				fmt.Println("deserialize Failure")
				continue
			}
			srv.recvChan <- msg
		}
	}()
}

func (srv *Server) ConnectWithBootnode() {
	srv.sendChan <- &UDPMessage{
		RemoteAddress: &srv.IP,
		ID:            srv.LocalNode.ID,
		data:          "hello bootnode",
		RemotePort:    8081,
		MessageClass:  Ping,
	}
}

func (srv *Server) processMessage() {

	recv := <-srv.recvChan

	switch recv.MessageClass {
	case Ping:
		fmt.Println("ping from node")
		peerNode := &Node{
			ID:   recv.ID,
			IP:   recv.LocalAddress,
			Port: int(recv.LocalPort),
		}
		srv.LocalNode.AddNode(peerNode)
	
	case Add: 
		
	}

}
