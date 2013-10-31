package main

import "fmt"

func main() {
	exit := make(chan int, 1)
	incoming_requests := make(chan Message, 200)
	fmt.Println("Starting!")
	go RunServer(exit, incoming_requests)
	go ManageRequests(incoming_requests)
	go fmt.Println("Server started. Press a key to exit.")
	fmt.Scanln()
	fmt.Println("Goodbye!")
	exit <- 1
	return
}
