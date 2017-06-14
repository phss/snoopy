package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	proxiedBaseUrl := os.Args[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxiedUrl := proxiedBaseUrl + r.URL.Path

		resp, _ := http.Get(proxiedUrl)
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "%s", body)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
