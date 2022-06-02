package router_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"test_go/internal/router"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tests = map[string]string{
	"Alloc":      "1111111",
	"TotalAlloc": "222222",
	"Sys":        "333333",
	"Mallocs":    "444444",
	"Frees":      "555555",
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	r := router.New()
	ts := httptest.NewServer(r)
	defer ts.Close()
	for k, v := range tests {
		resp, body := testRequest(t, ts, "POST", "/update/"+k+"/gauge/"+v)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		//fmt.Println(string(body))
		assert.Equal(t, k+"-"+"gauge"+"-"+v, body)
	}

}
