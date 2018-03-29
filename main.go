package main

import (
	"fmt"

	"github.com/alexchandra/message"
)

func main() {
	var (
		s = message.Message{
			Name: "Alex"
		}
	)

	fmt.Println(s)
}
