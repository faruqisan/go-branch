package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	// okEchoHandler is http handler which echoed the request with
	// 200 status code
	okEchoHandler = func(w http.ResponseWriter, r *http.Request) {
		for name, val := range r.Header {
			w.Header().Set(name, val[0])
		}

		w.WriteHeader(http.StatusOK)
		io.Copy(w, r.Body)
	}

	// helloJSONHandler is http handler which echoed the request headers
	// and give hello json in response body
	helloJSONHandler = func(w http.ResponseWriter, r *http.Request) {
		for name, val := range r.Header {
			w.Header().Set(name, val[0])
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"key":"value"}`)
	}

	// nopHandler is http handler that do nothing to the request
	nopHandler = func(w http.ResponseWriter, r *http.Request) {
	}
)

func TestHTTPClientDo(t *testing.T) {
	client := NewClient()
	testHTTPClientDo(t, client.Do)
}

func testHTTPClientDo(t *testing.T, doer func(*http.Request) (*http.Response, error)) {
	const (
		body = "my body"
	)

	server := httptest.NewServer(http.HandlerFunc(okEchoHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, strings.NewReader(body))
	require.NoError(t, err)

	resp, err := doer(req)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, body, string(respBody))
}

func TestHTTPClientDoJSON(t *testing.T) {
	client := NewClient()
	testHTTPClientDoJSON(t, client.DoJSON)
}

func testHTTPClientDoJSON(t *testing.T, doer func(*http.Request, interface{}) (*http.Response, error)) {

	testCases := []struct {
		name    string
		handler http.HandlerFunc
		body    string
		wantErr bool
	}{
		{
			name:    "valid json",
			handler: okEchoHandler,
			body:    `{"name":"you"}`,
			wantErr: false,
		},
		{
			name:    "invalid json",
			handler: okEchoHandler,
			body:    "name",
			wantErr: true,
		},
		{
			name:    "nop",
			handler: nopHandler,
			body:    "name",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(tc.handler))
			defer server.Close()

			req, err := http.NewRequest(http.MethodGet, server.URL, strings.NewReader(tc.body))
			require.NoError(t, err)

			data := make(map[string]interface{})

			_, err = doer(req, &data)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestHTTPClientGet(t *testing.T) {
	client := NewClient()

	server := httptest.NewServer(http.HandlerFunc(okEchoHandler))
	defer server.Close()

	// set the headers
	headers := map[string][]string{
		"K1": {"v1"},
		"K2": {"v2"},
	}

	// execute
	resp, err := client.Get(context.Background(), server.URL, headers)
	require.NoError(t, err)

	// check we got same headers as we sent
	for k := range headers {
		require.Equal(t, headers[k], resp.Header[k])
	}
}

func TestHTTPClientGetJSON(t *testing.T) {
	client := NewClient()

	server := httptest.NewServer(http.HandlerFunc(helloJSONHandler))
	defer server.Close()

	// set the headers
	headers := map[string][]string{
		"K1": {"v1"},
		"K2": {"v2"},
	}

	// execute
	data := make(map[string]string)
	resp, err := client.GetJSON(context.Background(), server.URL, headers, &data)
	require.NoError(t, err)

	// check we got same headers as we sent
	for k := range headers {
		require.Equal(t, headers[k], resp.Header[k])
	}

	// check we got the data we want
	require.Len(t, data, 1)
	require.Equal(t, data["key"], "value")
}

func TestHTTPClientPostWithBody(t *testing.T) {
	bodyObj := struct {
		Name string
		Age  int
	}{
		Name: "Tokopedia",
		Age:  9,
	}

	bodyBytes, err := json.Marshal(&bodyObj)
	require.Nil(t, err)

	testCases := []struct {
		name string
		body interface{}
	}{
		{
			name: "string",
			body: string(bodyBytes),
		},
		{
			name: "[]byte",
			body: bodyBytes,
		},

		{
			name: "nil",
			body: nil,
		},
		{
			name: "struct",
			body: &bodyObj,
		},
		{
			name: "io.Reader",
			body: bytes.NewReader(bodyBytes),
		},
	}

	client := NewClient()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(okEchoHandler))
			defer server.Close()

			require.NoError(t, err)

			resp, err := client.Post(context.Background(), server.URL, nil, tc.body)
			require.NoError(t, err)

			defer resp.Body.Close()

			respBody, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			if tc.body != nil {
				require.Equal(t, string(bodyBytes), string(respBody))
			} else {
				require.Equal(t, "", string(respBody))
			}
		})
	}
}

// TestHTTPClientTimeout test the client can do timeout properly
func TestHTTPClientTimeout(t *testing.T) {
	const (
		timeoutSec = 1
	)
	client := NewClient(WithHTTPTimeout(timeoutSec * time.Second))

	slowHandler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep((timeoutSec + 2) * time.Second)
		w.WriteHeader(http.StatusOK)
	}

	var server *httptest.Server
	startedCh := make(chan struct{})
	go func() {
		server = httptest.NewServer(http.HandlerFunc(slowHandler))
		startedCh <- struct{}{}
	}()

	<-startedCh
	defer server.Close()

	_, err := client.Get(context.Background(), server.URL, nil)
	require.Error(t, err)
}

func TestMaxConcurrentRequest(t *testing.T) {
	const (
		max = 10
	)

	client := NewClient(
		WithMaxConcurrentRequest(max),
		WithRequestVolumeThreshold(5),
	)

	// start the server
	slowHandler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}
	var server *httptest.Server
	startedCh := make(chan struct{})
	go func() {
		server = httptest.NewServer(http.HandlerFunc(slowHandler))
		startedCh <- struct{}{}
	}()

	// waiting for the server to be started
	<-startedCh

	// start request which exceeds the number of max requests
	var (
		numReq = max + 2
		wg     sync.WaitGroup
		errCh  = make(chan error, numReq)
	)

	wg.Add(numReq)

	for i := 0; i < numReq; i++ {
		go func() {
			defer wg.Done()
			_, err := client.Get(context.Background(), server.URL, nil)
			if err != nil {
				errCh <- err
			}
		}()
	}
	wg.Wait()

	require.Equal(t, numReq-max, len(errCh))
}
