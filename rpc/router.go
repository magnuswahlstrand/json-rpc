package rpc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type router struct {
	procedures map[string]ProcedureInf
}

func Router() *router {
	r := router{map[string]ProcedureInf{}}
	return &r
}

func (rou *router) AddProcedure(p ProcedureInf) *router {
	rou.procedures[p.Name()] = p
	return rou
}

func (rou *router) Handler(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte

	// TODO: replace with tee-reader?
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "aerror not supported: %s", err)
		return
	}

	var req map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		returnError(w, nil, -32700, "Parse error")
		return
	}
	id := req["id"]

	method, methodFound := req["method"].(string)
	if !methodFound {
		returnError(w, id, -32600, "Invalid Request")
		return
	}

	procedure, found := rou.procedures[method]
	if !found {
		returnError(w, id, -32601, "Method not found")
		return
	}

	out, err := procedure.serverCall(bodyBytes)
	if err != nil {
		fmt.Fprintf(w, "aerror not supported: %s", err)
		return
	}

	if id == nil {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	fmt.Fprintf(w, string(out))
}

func (rou *router) ListenAndServe(port string) error {
	return http.ListenAndServe(port, http.HandlerFunc(rou.Handler))
}

func returnError(w http.ResponseWriter, id interface{}, code int, message string) {
	v := Response[*struct{}]{
		JSONRPC: "2.0",
		Result:  nil,
		ID:      id,
		Error: &ErrorFoo{
			Code:    code,
			Message: message,
		},
	}
	_ = json.NewEncoder(w).Encode(v)
}
