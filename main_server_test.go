package main

import "testing"
import "bytes"
import "encoding/binary"
import "fmt"
import "net"

import "time"

func TestConnect(t *testing.T) {
	//exit := make(chan int, 1)
	//blah := make(chan Message, 200)
	//go RunServer(exit, blah)
	time.Sleep(1 * time.Second)
	num_conn := 1000
	conns := [1000]*net.UDPConn{}
	ra, err := net.ResolveUDPAddr("udp", "localhost:24816")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	for i := 0; i < num_conn; i++ {
		con, err := net.DialUDP("udp", nil, ra)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		conns[i] = con
	}
	fmt.Println("Connections Complete")
	original_message := "This is a test message!"
	message_bytes := []byte(original_message)
	var msg_len = make([]byte, 4)
	binary.LittleEndian.PutUint32(msg_len, uint32(len(message_bytes)))
	output_message := append(append([]byte{0, 0, 0, 0, 0}, msg_len...), message_bytes...)
	var buf [512]byte
	count := 0

	for i := 0; i < num_conn*50; i++ {
		var v = i % num_conn
		con := conns[v]
		_, err := con.Write(output_message)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		var c_len int32
		binary.Read(bytes.NewBuffer(buf[5:9]), binary.LittleEndian, &c_len)
		string_return := string(buf[9 : 9+c_len])
		if n != len(message_bytes)+9 || string_return != original_message {
			fmt.Println("Message length or content did not match!")
			t.FailNow()
		} else {
			count += 1
		}
		if count%num_conn == 0 {
			fmt.Println("Count: ", count)
		}
	}
	//exit <- 1
}

func BenchmarkEcho(b *testing.B) {
}
