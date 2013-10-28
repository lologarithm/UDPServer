package main

import (
	"fmt"
	"net"
	"os"
)

const (
	port string = ":24816"
)

type Server struct {
	conn        *net.UDPConn
	connections map[int]Client
}

type Client struct {
	userID   int
	buffer   []byte
	userAddr *net.UDPAddr
}

type Message struct {
	value string
}

func (s *Server) handleMessage() {
	buf := make([]byte, 512)

	n, addr, err := s.conn.ReadFromUDP(buf[0:])
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	fmt.Println("Got message from; ", addr)
	msg_obj := s.parseMessage(buf[0:n])
	fmt.Println("OBJ: ", msg_obj)
}

func (s *Server) parseMessage(raw_msg []byte) (m Message) {
	fmt.Println(raw_msg)
	m.value = string(raw_msg)
	fmt.Println("Got string: ", m.value)
	return
}

func (s *Server) sendMessage(msg string) {
	for _, c := range s.connections {
		fmt.Println(c)
		n, err := s.conn.WriteToUDP([]byte(msg), c.userAddr)
		fmt.Println(n, err)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error:%s", err.Error())
		os.Exit(1)
	}
}

func RunServer(exit chan int) {
	udpAddr, err := net.ResolveUDPAddr("udp4", port)
	checkError(err)

	var s Server
	s.connections = make(map[int]Client, 0)

	s.conn, err = net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {
		select {
		case <-exit:
			fmt.Println("Killing Socket Server")
			s.conn.Close()
			break
		default:
			fmt.Println("Reading.")
			s.handleMessage()
		}
	}
}
