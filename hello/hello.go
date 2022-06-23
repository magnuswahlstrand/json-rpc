package hello

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

func calculateGreeting(req Req) (Resp, error) {
	return Resp{Greeting: fmt.Sprintf("Hello, %s!", req.Name)}, nil
}

func subtract(req [2]int64) (int64, error) {
	return req[0] - req[1], nil
}

var CalculateGreeting = rpc.RemoteProcedure("hello", calculateGreeting)
var Subtract = rpc.RemoteProcedure("subtract", subtract)

var Router = rpc.Router().
	AddProcedure(CalculateGreeting).
	AddProcedure(Subtract)
