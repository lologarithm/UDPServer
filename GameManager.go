package main

import "fmt"
import "time"

func ManageRequests(incoming_requests chan Message) {
	for {
		select {
		case msg := <-incoming_requests:
			fmt.Println("MESSAGE:", msg)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
