package main

import (
	"github.com/magnuswahlstrand/json-rpc/examples/hello_world/world"
	"log"
)

func main() {
	if err := world.Router.ListenAndServe(":8081"); err != nil {
		log.Fatalln(err)
	}
}
