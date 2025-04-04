package main

import (
	"bufio"
	"fmt"
	"os"
)

func read() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("user> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return text
}

func eval(input string) string {
	return input
}

func print(output string) {
	fmt.Println(output)
}

func main() {
	for {
		input := read()
		evaled := eval(input)
		print(evaled)
	}
}
