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
	conn              *net.UDPConn
	players           map[int32]*Client
	connections       map[net.Addr]*Client
	outgoing_player   chan Message
	incoming_requests chan Message
	input_buffer      []byte
	encryption_key    *rsa.PrivateKey
}

func (s *Server) handleMessage() {
	n, addr, err := s.conn.ReadFromUDP(s.input_buffer)
	if err != nil {
		fmt.Println("ERROR: ", err)
		return
	}
	if n == 0 {
		// send exit signal to client
		close(s.connections[addr].incoming_bytes)
		delete(s.connections, addr) // Expire the client goroutine.
	}
	if _, ok := s.connections[addr]; !ok {
		s.connections[addr] = &Client{client_address: addr, incoming_bytes: make(chan []byte, 200)}
		go s.connections[addr].ProcessBytes(s.incoming_requests, s.outgoing_player)
	}
	s.connections[addr].incoming_bytes <- s.input_buffer[0:n]
}

func ParseFrame(raw_bytes []byte) *MessageFrame {
	if len(raw_bytes) > 9 {
		//fmt.Println("RAW BYTES:", raw_bytes)
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
			msg.destination = s.players[msg.destination.user.id]
		}
		if n, err := s.conn.WriteToUDP(msg.raw_bytes, msg.destination.client_address); err != nil {
			fmt.Println("Error: ", err, " Bytes Written: ", n)
		}
	}
}

type Client struct {
	buffer         []byte
	client_address *net.UDPAddr
	incoming_bytes chan []byte
	user           User
}

func (client *Client) ProcessBytes(to_client chan Message, outgoing_msg chan Message) {
	for {
		dem_bytes, ok := <-client.incoming_bytes
		if !ok {
			break
		}
		client.buffer = append(client.buffer, dem_bytes...)
		msg_frame := ParseFrame(client.buffer)
		if msg_frame != nil && int(msg_frame.frame_length+msg_frame.content_length) >= len(client.buffer) {
			msg_obj := client.parseMessage(msg_frame)
			if msg_obj.frame.message_type == 0 {
				msg_obj.destination = client
				//fmt.Println("Sending message out.")
				outgoing_msg <- msg_obj
			} else {
				to_client <- msg_obj
			}
		}
	}
}

func (client *Client) parseMessage(mf *MessageFrame) (m Message) {
	m.raw_bytes = client.buffer[0 : mf.frame_length+mf.content_length]
	m.frame = mf
	client.buffer = client.buffer[mf.frame_length+mf.content_length:]
	return
}

type Message struct {
	raw_bytes   []byte
	frame       *MessageFrame
	destination *Client
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

func RunServer(exit chan int, requests chan Message) {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	checkError(err)
	fmt.Println("Now listening on port", port)

	var s Server
	s.connections = make(map[net.Addr]*Client, 0)
	s.input_buffer = make([]byte, 512)
	s.outgoing_player = make(chan Message, 255)
	s.incoming_requests = requests
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
			fmt.Println("Looking for new messages")
			s.handleMessage()
		}
	}
}
