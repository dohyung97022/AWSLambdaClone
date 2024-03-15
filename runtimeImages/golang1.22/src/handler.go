package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

func handler(params url.Values, w *http.ResponseWriter) {
	json, _ := json.Marshal(map[string]any{"message": "Hello World", "params": params})
	(*w).Write(json)
	(*w).WriteHeader(200)
}
