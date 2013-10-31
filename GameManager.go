package main

import "fmt"

func ManageRequests(incoming_requests chan Message) {
	for {
		select {
		case msg := <-incoming_requests:
			fmt.Println("MESSAGE:", msg)
		default:

		}
	}
}
