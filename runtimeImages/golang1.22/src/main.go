package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/*", func(w http.ResponseWriter, req *http.Request) {
		handler(req.URL.Query(), &w)
	})

	if err := http.ListenAndServe(":443", nil); err != nil {
		fmt.Println(err)
		return
	}
}
