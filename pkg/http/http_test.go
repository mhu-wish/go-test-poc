package http

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeHTTPCall(t *testing.T) {
	testTable := []struct {
		name             string
		server           *httptest.Server
		expectedResponse *Response
		expectedErr      bool
	}{
		{
			name: "happy-server-response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"id": 1, "name": "kyle", "description": "novice gopher"}`))
			})),
			expectedResponse: &Response{
				ID:          1,
				Name:        "kyle",
				Description: "novice gopher",
			},
			expectedErr: false,
		},
		{
			name: "invalid-json-response",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`123`))
			})),
			expectedResponse: nil,
			expectedErr: true,
		},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.server.Close()
			resp, err := MakeHTTPCall(tc.server.URL)
			require.Equal(t, tc.expectedResponse, resp)
			require.Equal(t, tc.expectedErr, err != nil)
		})
	}
}
