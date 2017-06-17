package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
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

func printResponse(resp *http.Response, body string) {
	color.Blue("Response")
	fmt.Printf("%s: %s\n", color.CyanString("Status"), color.YellowString(resp.Status))
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", color.CyanString(name), color.YellowString(value))
		}
	}
	fmt.Printf("%s:\n", color.CyanString("Body"))
	fmt.Println(body)
}

func makeProxyRequest(proxiedUrl string, r *http.Request) (*http.Response, string) {
	client := http.Client{}
	proxyRequest, _ := http.NewRequest(r.Method, proxiedUrl, r.Body)
	resp, _ := client.Do(proxyRequest)
	body, _ := ioutil.ReadAll(resp.Body)
	return resp, string(body)
}

func returnProxyResponse(resp *http.Response, body string, w http.ResponseWriter) {
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	defer resp.Body.Close()
	fmt.Fprintf(w, "%s", body)
}

func proxy(port int, proxiedBaseUrl string) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxiedUrl := proxiedBaseUrl + r.URL.Path
		printRequest(proxiedUrl, r)
		resp, body := makeProxyRequest(proxiedUrl, r)
		printResponse(resp, body)
		returnProxyResponse(resp, body, w)
	})

	fmt.Printf("Proxying for %s on http://localhost:%d \n", proxiedBaseUrl, port)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	log.Fatal(server.ListenAndServe())
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		proxy(8080, os.Args[1])
		wg.Done()
	}()
	go func() {
		proxy(8081, os.Args[2])
		wg.Done()
	}()

	wg.Wait()
}
