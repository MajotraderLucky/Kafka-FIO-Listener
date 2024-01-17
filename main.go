package main

import (
	"fmt"
	"time"
)

func main() {
	// Wait for other services to be ready
	time.Sleep(time.Second * 10)

	fmt.Println("Hello, world!")
}
