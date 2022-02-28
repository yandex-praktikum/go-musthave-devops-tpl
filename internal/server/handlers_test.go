package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/itd27m01/go-metrics-service/internal/pkg/metrics"
	"github.com/itd27m01/go-metrics-service/internal/repository"
	"github.com/itd27m01/go-metrics-service/internal/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const metricsHTML = `<!DOCTYPE html>
<html lang="en">
<body>
<table>
    <tr>
        <th>Type</th>
        <th>Name</th>
        <th>Value</th>
    </tr>
    <tr>
        <td style='text-align:center; vertical-align:middle'>gauge</td>
        <td style='text-align:center; vertical-align:middle'>test1</td>
        <td style='text-align:center; vertical-align:middle'>100</td>
    </tr>
    <tr>
        <td style='text-align:center; vertical-align:middle'>counter</td>
        <td style='text-align:center; vertical-align:middle'>test2</td>
        <td style='text-align:center; vertical-align:middle'>100</td>
    </tr>
    <tr>
        <td style='text-align:center; vertical-align:middle'>gauge</td>
        <td style='text-align:center; vertical-align:middle'>testSetGet134</td>
        <td style='text-align:center; vertical-align:middle'>96969.519</td>
    </tr>
    <tr>
        <td style='text-align:center; vertical-align:middle'>gauge</td>
        <td style='text-align:center; vertical-align:middle'>testSetGet135</td>
        <td style='text-align:center; vertical-align:middle'>156519.255</td>
    </tr>
    </table>
</body>
</html>`

type want struct {
	code int
	data string
}

type test struct {
	name   string
	method string
	metric string
	want   want
}

type testJSON struct {
	name   string
	method string
	url    string
	metric *metrics.Metric
	want   want
}

var gaugeValue metrics.Gauge = 96969.519

var testsJSON = []testJSON{
	{
		name:   "Post JSON metric",
		method: http.MethodPost,
		url:    "/update/",
		metric: &metrics.Metric{
			ID:    "Alloc",
			MType: metrics.GaugeMetricTypeName,
			Value: &gaugeValue,
		},
		want: want{
			code: http.StatusOK,
			data: "",
		},
	},
	{
		name:   "Get JSON metric",
		method: http.MethodPost,
		url:    "/value/",
		metric: &metrics.Metric{
			ID:    "Alloc",
			MType: metrics.GaugeMetricTypeName,
		},
		want: want{
			code: http.StatusOK,
			data: "{\"id\":\"Alloc\",\"type\":\"gauge\",\"value\":96969.519}\n",
		},
	},
}

var tests = []test{
	{
		name:   "OK gauge update",
		metric: "/update/gauge/test1/100.000000",
		method: http.MethodPost,
		want: want{
			code: http.StatusOK,
		},
	},
	{
		name:   "OK counter update",
		metric: "/update/counter/test2/100",
		method: http.MethodPost,
		want: want{
			code: http.StatusOK,
		},
	},
	{
		name:   "Test gauge post 1",
		metric: "/update/gauge/testSetGet134/96969.519",
		method: http.MethodPost,
		want: want{
			code: http.StatusOK,
		},
	},
	{
		name:   "Test gauge post 2",
		metric: "/update/gauge/testSetGet135/156519.255",
		method: http.MethodPost,
		want: want{
			code: http.StatusOK,
		},
	},
	{
		name:   "Test gauge get 1",
		metric: "/value/gauge/testSetGet134",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: "96969.519",
		},
	},
	{
		name:   "Test gauge get 2",
		metric: "/value/gauge/testSetGet135",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: "156519.255",
		},
	},
	{
		name:   "BAD gauge update",
		metric: "/update/gauge/test/none",
		method: http.MethodPost,
		want: want{
			code: http.StatusBadRequest,
		},
	},
	{
		name:   "BAD counter update",
		metric: "/update/counter/test/none",
		method: http.MethodPost,
		want: want{
			code: http.StatusBadRequest,
		},
	},
	{
		name:   "NotFound gauge update",
		metric: "/update/gauge/",
		method: http.MethodPost,
		want: want{
			code: http.StatusNotFound,
		},
	},
	{
		name:   "NotFound counter update",
		metric: "/update/counter/",
		method: http.MethodPost,
		want: want{
			code: http.StatusNotFound,
		},
	},
	{
		name:   "NotImplemented update",
		metric: "/update/unknown/test/1001",
		method: http.MethodPost,
		want: want{
			code: http.StatusNotImplemented,
		},
	},
	{
		name:   "Get all metrics",
		metric: "/",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: metricsHTML,
		},
	},
	{
		name:   "Get gauge metric",
		metric: "/value/gauge/test1",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: "100",
		},
	},
	{
		name:   "Get counter metric",
		metric: "/value/counter/test2",
		method: http.MethodGet,
		want: want{
			code: http.StatusOK,
			data: "100",
		},
	},
	{
		name:   "Get unknown metric",
		metric: "/value/counter/unknown",
		method: http.MethodGet,
		want: want{
			code: http.StatusNotFound,
			data: "Metric not found: unknown\n",
		},
	},
}

func TestRouter(t *testing.T) {
	mux := chi.NewRouter()
	server.RegisterHandlers(mux, repository.NewInMemoryStore())
	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRequest(t, ts, tt)
		})
	}

	for _, tt := range testsJSON {
		t.Run(tt.name, func(t *testing.T) {
			testJSONRequest(t, ts, tt)
		})
	}
}

func testJSONRequest(t *testing.T, ts *httptest.Server, testData testJSON) {
	body, _ := testData.metric.EncodeMetric()
	req, err := http.NewRequest(testData.method, ts.URL+testData.url, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, testData.want.code, resp.StatusCode)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	if string(respBody) != "" {
		assert.JSONEq(t, testData.want.data, string(respBody))
	} else {
		assert.Equal(t, testData.want.data, string(respBody))
	}

	require.NoError(t, err)

	err = resp.Body.Close()
	if err != nil {
		return
	}
}

func testRequest(t *testing.T, ts *httptest.Server, testData test) {
	req, err := http.NewRequest(testData.method, ts.URL+testData.metric, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.Equal(t, testData.want.code, resp.StatusCode)
	require.NoError(t, err)
	defer resp.Body.Close()

	if testData.method == http.MethodGet {
		respBody, err := ioutil.ReadAll(resp.Body)
		assert.Equal(t, testData.want.data, string(respBody))
		require.NoError(t, err)
	}
}
