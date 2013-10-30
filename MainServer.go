package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

const (
	port string = ":24816"
)

type Server struct {
	conn            *net.UDPConn
	players         map[int]Client
	connections     map[net.Addr]Client
	outgoing_player chan Message
	input_buffer    []byte
	encryption_key  *rsa.PrivateKey
}

func (s *Server) handleMessage() {
	n, addr, err := s.conn.ReadFromUDP(s.input_buffer)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}

	if _, ok := s.connections[addr]; !ok {
		s.connections[addr] = Client{client_address: addr}
	}
	client := s.connections[addr]
	client.buffer = append(client.buffer, s.input_buffer[0:n]...)
	msg_frame := ParseFrame(client.buffer)
	if msg_frame != nil && int(msg_frame.frame_length+msg_frame.content_length) >= len(client.buffer) {
		msg_obj := s.parseMessage(&client, msg_frame)
		fmt.Println("Message: ", msg_obj.raw_bytes)
		fmt.Println("    Message Type:\t", msg_obj.frame.message_type)
		fmt.Println("    Content Len:\t", msg_obj.frame.content_length)
		fmt.Println("    Content: \t\t", msg_obj.Content())
		msg_obj.destination = client
		if msg_obj.frame.message_type == 0 {
			s.outgoing_player <- msg_obj
		}
	}
}

func (s *Server) parseMessage(client *Client, mf *MessageFrame) (m Message) {
	m.raw_bytes = client.buffer[0 : mf.frame_length+mf.content_length]
	m.frame = mf
	client.buffer = client.buffer[mf.frame_length+mf.content_length:]
	return
}

func ParseFrame(raw_bytes []byte) *MessageFrame {
	if len(raw_bytes) > 9 {
		fmt.Println("RAW BYTES:", raw_bytes)
		mf := new(MessageFrame)
		mf.message_type = raw_bytes[0]
		var v int32
		binary.Read(bytes.NewBuffer(raw_bytes[1:5]), binary.LittleEndian, &v)
		mf.from_user = v
		binary.Read(bytes.NewBuffer(raw_bytes[5:9]), binary.LittleEndian, &v)
		mf.content_length = v
		mf.frame_length = 9
		return mf
	}

	return nil
}

func (s *Server) sendMessages() {
	for {
		msg := <-s.outgoing_player
		if msg.destination.client_address == nil {
			msg.destination = s.players[msg.destination.player.id]
		}
		if n, err := s.conn.WriteToUDP(msg.raw_bytes, msg.destination.client_address); err != nil {
			fmt.Println("Error: ", err, " Bytes Written: ", n)
		}
	}
}

type Client struct {
	buffer         []byte
	client_address *net.UDPAddr
	player         Player
}

type Message struct {
	raw_bytes   []byte
	frame       *MessageFrame
	destination Client
}

func (m *Message) Content() []byte {
	return m.raw_bytes[m.frame.frame_length : m.frame.frame_length+m.frame.content_length]
}

type MessageFrame struct {
	message_type   byte
	from_user      int32
	frame_length   int32
	content_length int32
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
	s.connections = make(map[net.Addr]Client, 0)
	s.input_buffer = make([]byte, 512)
	s.outgoing_player = make(chan Message, 255)
	s.conn, err = net.ListenUDP("udp", udpAddr)
	checkError(err)

	go s.sendMessages()

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
