package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getHandler(sf sunsetFinder) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		location := req.URL.Query().Get("location")

		if location == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("location must be set"))
			return
		}

		result, err := sf.Query(location)

		if err == errNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Location '%s' not found", location)))
		} else if err != nil {
			w.WriteHeader(http.StatusBadGateway)
		}

		// for example purposes only, just assume
		// this won't fail
		b, _ := json.Marshal(result)

		w.Write(b)
	})
}
