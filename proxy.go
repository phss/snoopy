package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	port := 8080
	proxiedBaseUrl := os.Args[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxiedUrl := proxiedBaseUrl + r.URL.Path

		fmt.Println("---")
		color.Blue("Request")
		color.Green("%s %s -> %s", r.Method, r.URL.Path, proxiedUrl)
		for name, values := range r.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", color.CyanString(name), color.YellowString(value))
			}
		}
		fmt.Println("")

		resp, _ := http.Get(proxiedUrl)
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Fprintf(w, "%s", body)

		color.Blue("Response")
		for name, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", color.CyanString(name), color.YellowString(value))
			}
		}
		//fmt.Println("")
		//fmt.Println(string(body))
	})

	fmt.Printf("Proxying for %s on http://localhost:%d \n", proxiedBaseUrl, port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
