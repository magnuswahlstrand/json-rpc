package main

import (
	"fmt"
	"github.com/magnuswahlstrand/json-rpc/examples/hello_world/world"
	"log"
)

const url = "http://localhost:8081"

var greet = world.Greet.Bind(url)

func main() {
	resp, err := greet(world.Req{Name: "Magnus"})
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("received response: %q\n", resp.Greeting)
}
