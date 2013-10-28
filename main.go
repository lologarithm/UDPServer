package main

import "fmt"

func main() {
	exit := make(chan int, 1)
	fmt.Println("Starting!")
	go RunServer(exit)
	fmt.Println("Server started. Press a key to exit.")
	fmt.Scanln()
	fmt.Println("Goodbye!")
	exit <- 1
	return
}
