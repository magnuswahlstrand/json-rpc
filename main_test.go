package main

import (
	"github.com/magnuswahlstrand/json-rpc/rpc"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

func TestRPC(t *testing.T) {
	tcs := []struct {
		procedure rpc.ProcedureInf
		name      string
		in        string
		expected  string
	}{
		{
			procedure: rpc.RemoteProcedure("subtract", subtract),
			name:      "rpc call with positional parameters 1",
			in:        `{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}`,
			expected:  `{"jsonrpc":"2.0","result":19,"id":1}`,
		},
		{
			procedure: rpc.RemoteProcedure("subtract", subtract),
			name:      "rpc call with positional parameters 2",
			in:        `{"jsonrpc": "2.0", "method": "subtract", "params": [23, 42], "id": 2}`,
			expected:  `{"jsonrpc":"2.0","result":-19,"id":2}`,
		},
		{
			procedure: rpc.RemoteProcedure("subtract", namedSubtract),
			name:      "rpc call with named parameters 1",
			in:        `{"jsonrpc": "2.0", "method": "subtract", "params": {"subtrahend": 23, "minuend": 42}, "id": 3}`,
			expected:  `{"jsonrpc": "2.0","result": 19, "id": 3}`,
		},
		{
			procedure: rpc.RemoteProcedure("subtract", namedSubtract),
			name:      "rpc call with named parameters 2",
			in:        `{"jsonrpc": "2.0", "method": "subtract", "params": {"minuend": 42, "subtrahend": 23}, "id": 4}`,
			expected:  `{"jsonrpc": "2.0", "result":19, "id":4}`,
		},
		{
			procedure: rpc.RemoteProcedure("subtract", namedSubtract),
			name:      "rpc call of non-existent method",
			in:        `{"jsonrpc": "2.0", "method": "foobar", "id": "1"}`,
			expected:  `{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "1"}`,
		},
		{
			procedure: rpc.RemoteProcedure("subtract", namedSubtract),
			name:      "rpc call with invalid JSON",
			in:        `{"jsonrpc": "2.0", "method": "foobar, "params": "bar", "baz]`,
			expected:  `{"jsonrpc": "2.0", "error": {"code": -32700, "message": "Parse error"}, "id": null}`,
		},
		{
			procedure: rpc.RemoteProcedure("subtract", namedSubtract),
			name:      "rpc call with invalid Request object",
			in:        `{"jsonrpc": "2.0", "method": 1, "params": "bar"}`,
			expected:  `{"jsonrpc": "2.0", "error": {"code": -32600, "message": "Invalid Request"}, "id": null}`,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.in))
			w := httptest.NewRecorder()
			handler := rpc.Router().AddProcedure(tc.procedure).Handler

			// Act
			handler(w, req)
			actual, err := ioutil.ReadAll(w.Result().Body)

			// Assert
			require.NoError(t, err)
			require.JSONEq(t, tc.expected, string(actual))
		})
	}
}

func TestRPCNotifications(t *testing.T) {
	tcs := []struct {
		procedure rpc.ProcedureInf
		name      string
		in        string
	}{
		{
			procedure: rpc.RemoteProcedure("update", subtract),
			name:      "notification 1",
			in:        `{"jsonrpc": "2.0", "method": "update", "params": [1,2,3,4,5]}`,
		},
		{
			procedure: rpc.RemoteProcedure("foobar", func(t struct{}) (string, error) { return "", nil }),
			name:      "notification 2",
			in:        `{"jsonrpc": "2.0", "method": "foobar"}`,
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.in))
			w := httptest.NewRecorder()
			handler := rpc.Router().AddProcedure(tc.procedure).Handler

			// Act
			handler(w, req)
			actual, err := ioutil.ReadAll(w.Result().Body)

			// Assert
			require.NoError(t, err)
			require.Empty(t, string(actual))
		})
	}
}
