package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Request[T any] struct {
	JSONRPC string       `json:"jsonrpc"`
	Method  string       `json:"method"`
	Params  *T           `json:"params"`
	ID      *interface{} `json:"id"`
}

type Response[V any] struct {
	JSONRPC string       `json:"jsonrpc"`
	Result  V            `json:"result"`
	ID      *interface{} `json:"id"`
}

func subtract(in [2]int) (int, error) {
	return in[0] - in[1], nil
}

type namedSubtractRequest struct {
	Minuend    int `json:"minuend"`
	Subtrahend int `json:"subtrahend"`
}

func namedSubtract(in namedSubtractRequest) (int, error) {
	return in.Minuend - in.Subtrahend, nil
}

func aNotifiedFunction(in []int) (int, error) {
	return 0, nil
}

func NewHandler[T, V any](fn func(T) (V, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request[T]
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			panic("not implemented")
		}

		// Safeguards
		if req.Params == nil {
			panic("not supported")
		}

		resp, err := fn(*req.Params)
		if err != nil {
			panic("not implemented")
		}

		if req.ID == nil {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		rpcResp := Response[V]{
			"2.0",
			resp,
			req.ID,
		}

		if err := json.NewEncoder(w).Encode(rpcResp); err != nil {
			panic("not implemented")
		}
	}
}

func TestRPC(t *testing.T) {
	tcs := []struct {
		handler  http.HandlerFunc
		name     string
		in       string
		expected string
	}{
		{
			handler:  NewHandler(subtract),
			name:     "rpc call with positional parameters 1",
			in:       `{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}`,
			expected: `{"jsonrpc":"2.0","result":19,"id":1}`,
		},
		{
			handler:  NewHandler(subtract),
			name:     "rpc call with positional parameters 2",
			in:       `{"jsonrpc": "2.0", "method": "subtract", "params": [23, 42], "id": 2}`,
			expected: `{"jsonrpc":"2.0","result":-19,"id":2}`,
		},
		{
			handler:  NewHandler(namedSubtract),
			name:     "rpc call with named parameters 1",
			in:       `{"jsonrpc": "2.0", "method": "subtract", "params": {"subtrahend": 23, "minuend": 42}, "id": 3}`,
			expected: `{"jsonrpc":"2.0","result":19,"id":3}`,
		},
		{
			handler:  NewHandler(namedSubtract),
			name:     "rpc call with named parameters 2",
			in:       `{"jsonrpc": "2.0", "method": "subtract", "params": {"minuend": 42, "subtrahend": 23}, "id": 4}`,
			expected: `{"jsonrpc":"2.0","result":19,"id":4}`,
		},
		{
			handler:  NewHandler(namedSubtract),
			name:     "rpc call of non-existent method",
			in:       `{"jsonrpc": "2.0", "method": "foobar", "id": "1"}`,
			expected: `{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "1"}`,
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", strings.NewReader(tc.in))
			w := httptest.NewRecorder()
			tc.handler(w, req)
			data, err := ioutil.ReadAll(w.Result().Body)

			actual := strings.TrimSpace(string(data))
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestRPCNotifications(t *testing.T) {
	tcs := []struct {
		handler http.HandlerFunc
		name    string
		in      string
	}{
		{
			handler: NewHandler(func(t []int) (int, error) { return 0, nil }),
			name:    "notification 1",
			in:      `{"jsonrpc": "2.0", "method": "update", "params": [1,2,3,4,5]}`,
		},
		//{
		//	handler: NewHandler(func(t *string) (int, error) { return 0, nil }),
		//	name:    "notification 2",
		//	in:      `{"jsonrpc": "2.0", "method": "foobar"}`,
		//},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/upper?word=abc", strings.NewReader(tc.in))
			w := httptest.NewRecorder()
			tc.handler(w, req)
			data, err := ioutil.ReadAll(w.Result().Body)
			require.NoError(t, err)

			actual := strings.TrimSpace(string(data))
			require.Empty(t, actual)
		})
	}
}
