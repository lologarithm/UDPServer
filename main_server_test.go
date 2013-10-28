package main

import "testing"
import "fmt"
import "net"
import "time"

func TestConnect(t *testing.T) {
	exit := make(chan int, 1)
	go RunServer(exit)
	fmt.Println("Test starting")
	time.Sleep(time.Duration(1) * time.Second)
	ra, err := net.ResolveUDPAddr("udp4", "localhost:24816")
	con, err := net.DialUDP("udp4", nil, ra)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println("Sending test message.")
	_, err = con.Write([]byte("This is a test."))
	fmt.Println("Message sent.")
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	time.Sleep(time.Duration(1) * time.Second)
	exit <- 1
	time.Sleep(time.Duration(1) * time.Second)
}
