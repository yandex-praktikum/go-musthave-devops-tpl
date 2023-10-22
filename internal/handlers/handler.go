package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func StatusHandler(rw http.ResponseWriter, r *http.Request) {
	// ParseForm will parse query string values and make r.Form available
	//r.ParseForm()

	// r.Form is map of query string parameters
	// its' type is url.Values, which in turn is a map[string][]string
	queryMap := r.URL.Path

	switch r.Method {
	case http.MethodGet:
		// Handle GET requests
		rw.Header().Set("Content-Type", "plaint/text")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("Query string values: %s", queryMap)))
		return
	case http.MethodPost:
		// Handle POST requests
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			// Error occurred while parsing request body
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(fmt.Sprintf("Query string values: %s\nBody posted: %s", queryMap, body)))
		return
	}

	// Other HTTP methods (eg PUT, PATCH, etc) are not handled by the above
	// so inform the client with appropriate status code

	rw.WriteHeader(http.StatusMethodNotAllowed)
}
