package world

import (
	"fmt"
	"github.com/magnuswahlstrand/json-rpc/rpc"
)

type Req struct {
	Name string `json:"name"`
}

type Resp struct {
	Greeting string `json:"greeting"`
}

func helloWorld(req Req) (Resp, error) {
	return Resp{Greeting: fmt.Sprintf("Hello, %s!", req.Name)}, nil
}

var Greet = rpc.RemoteProcedure("hello", helloWorld)

var Router = rpc.Router().AddProcedure(Greet)
