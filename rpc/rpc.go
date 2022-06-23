package rpc

// TODO: Merge request types
type FooRequest struct {
	Method string      `json:"method"`
	ID     interface{} `json:"id"`
}

type Request[T any] struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  T           `json:"params"`
	ID      interface{} `json:"id"`
}

type Response[V any] struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  V           `json:"result,omitempty"`
	ID      interface{} `json:"id"`
	Error   *ErrorFoo   `json:"error,omitempty"`
}

type ErrorFoo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

//func RemoteProcedure[T, V any](fn func(T) (V, error)) func([]byte) ([]byte, error) {
//	return func(in []byte) ([]byte, error) {
//		var req Request[T]
//		if err := json.Unmarshal(in, &req); err != nil {
//			return nil, fmt.Errorf("error not supported: %s", err)
//		}
//
//		// Safeguards
//		if req.Params == nil {
//			return nil, fmt.Errorf("error not supported: parameter missing")
//		}
//
//		resp, err := fn(*req.Params)
//		if err != nil {
//			return nil, fmt.Errorf("error not supported: %s", err)
//		}
//
//		//if req.ID == nil {
//		//	w.WriteHeader(http.StatusAccepted)
//		//	return
//		//}
//
//		rpcResp := Response[V]{"2.0", resp, req.ID}
//
//		out, err := json.Marshal(rpcResp)
//		if err != nil {
//			return nil, fmt.Errorf("error not supported: %s", err)
//		}
//		return out, nil
//	}
//}
