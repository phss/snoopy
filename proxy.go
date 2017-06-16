package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func printRequest(proxiedUrl string, r *http.Request) {
	fmt.Println("---")
	color.Blue("Request")
	color.Green("%s %s -> %s", r.Method, r.URL.Path, proxiedUrl)
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", color.CyanString(name), color.YellowString(value))
		}
	}
	fmt.Printf("%s:\n", color.CyanString("Body"))
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	fmt.Println("")
}

func printResponse(resp *http.Response) {
	color.Blue("Response")
	fmt.Printf("%s: %s\n", color.CyanString("Status"), color.YellowString(resp.Status))
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", color.CyanString(name), color.YellowString(value))
		}
	}
	fmt.Printf("%s:\n", color.CyanString("Body"))
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func makeProxyRequest(proxiedUrl string, r *http.Request) *http.Response {
	client := http.Client{}
	proxyRequest, _ := http.NewRequest(r.Method, proxiedUrl, r.Body)
	resp, _ := client.Do(proxyRequest)
	return resp
}

func returnProxyResponse(resp *http.Response, w http.ResponseWriter) {
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Fprintf(w, "%s", body)
}

func main() {
	port := 8080
	proxiedBaseUrl := os.Args[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxiedUrl := proxiedBaseUrl + r.URL.Path
		printRequest(proxiedUrl, r)
		resp := makeProxyRequest(proxiedUrl, r)
		printResponse(resp)
		returnProxyResponse(resp, w)
	})

	fmt.Printf("Proxying for %s on http://localhost:%d \n", proxiedBaseUrl, port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
