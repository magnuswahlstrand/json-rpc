package main

import (
	"github.com/magnuswahlstrand/json-rpc/hello"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello.Router.Handler)
	http.ListenAndServe(":8080", mux)
}
