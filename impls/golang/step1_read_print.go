package main

import (
	"bufio"
	"fmt"
	"mal/reader"
	"mal/types"
	"os"
)

func read() string {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("user> ")
	text, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return text
}

func eval(input string) (types.MalType, bool) {
	r := reader.NewReader(input)
	return reader.GetAST(r)
}

func print(output types.MalType) {
	fmt.Println(output.GetStr())
}

func main() {
	for {
		input := read()
		evaled, ok := eval(input)
		if ok {
			print(evaled)
		}
	}
}
