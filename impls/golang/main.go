package main

import (
	"bufio"
	"fmt"
	"mal/core"
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

func eval(input string, e *env.Env) (types.MalType, bool) {
	r := reader.NewReader(input)
	ast, ok := parser.Parse(r)
	if !ok {
		return nil, false
	}
	return evaluator.Eval(ast, e)
}

func print(output types.MalType) {
	fmt.Println(output.GetStr(true))
}

func main() {
	env := env.NewEnv(nil)
	core.AddCoreToEnv(env)
	for {
		input := read()
		evaled, ok := eval(input, env)
		if ok {
			print(evaled)
		}
	}
}
