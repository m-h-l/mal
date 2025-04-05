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
		Fns: []env.Fn{
			env.NewBuiltinFn("+", func(args ...types.MalType) []types.MalType {
				if len(args) >= 2 && args[0].GetTypeId() == types.Number && args[1].GetTypeId() == types.Number {
					r := args[0].(*types.MalNumber).Add(*(args[1].(*types.MalNumber)))
					return append([]types.MalType{&r}, args[2:]...)
				}
				return nil
			}),
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
