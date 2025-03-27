package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("user> ")
		text, _ := reader.ReadString('\n')
		fmt.Println(text)
	}
}
