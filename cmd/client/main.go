package main

import (
	"fmt"
	"github.com/magnuswahlstrand/json-rpc/hello"
	"log"
)

func main() {
	resp, err := hello.CalculateGreeting.Call(hello.Req{Name: "Magnus"})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("success!", resp.Greeting)

	resp2, err := hello.Subtract.Call([2]int64{57, 15})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("success!", resp2)
}
