package main

import (
	"bufio"
	"fmt"
	"mal/env"
	"mal/evaluator"
	"mal/parser"
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
	ast, ok := parser.Parse(r)
	if !ok {
		return nil, false
	}
	env := env.Env{
		Fns: []types.Fn{
			types.Plus(),
			types.Multiply(),
			types.Subtract(),
			types.Divide(),
		},
	}
	return evaluator.Eval(ast, env), true
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
