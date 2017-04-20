package main

import (
	"fmt"
	"time"
)

func main() {
	go doSomething()
	doSometihngElse()
}

func doSomething() {
  for {}
  for i := 0; i < 5; i++ {
    fmt.Printf("something %v\n", i)
		time.Sleep(time.Millisecond * 200)
	}
}

func doSometihngElse() {
	for i := 0; i < 5; i++ {
		fmt.Printf("else %v\n", i)
		time.Sleep(time.Millisecond * 200)
	}
}
