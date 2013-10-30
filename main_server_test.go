package main

import "testing"
import "bytes"
import "encoding/binary"
import "fmt"
import "net"

func TestConnect(t *testing.T) {
	exit := make(chan int, 1)
	go RunServer(exit)
	ra, err := net.ResolveUDPAddr("udp4", "localhost:24816")
	con, err := net.DialUDP("udp4", nil, ra)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	original_message := "This is a test message!"
	message_bytes := []byte(original_message)
	var msg_len = make([]byte, 4)
	binary.LittleEndian.PutUint32(msg_len, uint32(len(message_bytes)))
	_, err = con.Write(append(append([]byte{0, 0, 0, 0, 0}, msg_len...), message_bytes...))
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	var buf [512]byte
	n, err := con.Read(buf[0:])
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	var c_len int32
	binary.Read(bytes.NewBuffer(buf[5:9]), binary.LittleEndian, &c_len)
	string_return := string(buf[9 : 9+c_len])
	fmt.Println("Total Bytes: ", n, "Message: ", string_return)
	if n != len(message_bytes)+9 || string_return != original_message {
		fmt.Println("Message length or content did not match!")
		t.FailNow()
	}
	exit <- 1
}
