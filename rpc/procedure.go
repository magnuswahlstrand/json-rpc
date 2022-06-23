package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func RemoteProcedure[T any, V any](name string, fn func(T) (V, error)) *procedure[T, V] {
	return &procedure[T, V]{name: name, fn: fn}
}

type procedure[T, V any] struct {
	name string
	fn   func(T) (V, error)
}

func (p *procedure[T, V]) Name() string {
	return p.name
}

func (p *procedure[T, V]) serverCall(in []byte) ([]byte, error) {
	fmt.Println(string(in))
	var req Request[T]
	if err := json.Unmarshal(in, &req); err != nil {
		return nil, fmt.Errorf("error not supported: %s", err)
	}

	resp, err := p.fn(req.Params)
	if err != nil {
		return nil, fmt.Errorf("error not supported: %s", err)
	}
	rpcResp := Response[V]{"2.0", resp, req.ID, nil}

	out, err := json.Marshal(rpcResp)
	if err != nil {
		return nil, fmt.Errorf("error not supported: %s", err)
	}
	return out, nil
}

func (p *procedure[T, V]) Call(iu string, input T) (V, error) {
	u, err := url.Parse(iu)
	if err != nil {
		return *new(V), err
	}

	id := "123"
	req := Request[T]{
		JSONRPC: "2.0",
		Method:  p.name,
		Params:  input,
		ID:      id,
	}

	b, err := json.Marshal(req)
	if err != nil {
		return *new(V), err
	}

	resp, err := http.DefaultClient.Post(u.String(), "application/json", bytes.NewReader(b))
	if err != nil {
		return *new(V), err
	}

	var out Response[V]
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return *new(V), err
	}
	return out.Result, nil
}

func (p *procedure[T, V]) Bind(url string) func(T) (V, error) {
	return func(in T) (V, error) {
		return p.Call(url, in)
	}
}

type ProcedureInf interface {
	Name() string
	serverCall([]byte) ([]byte, error)
}
