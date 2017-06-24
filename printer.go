package main

import (
	"fmt"
	"github.com/fatih/color"
	"net/http"
)

func printRequest(proxyRequest ProxyRequest, config Config) {
	fmt.Println("---")
	color.Blue("Request")
	color.Green("%s %s -> %s", proxyRequest.Method, proxyRequest.Path, proxyRequest.ProxiedUrl())
	for _, header := range proxyRequest.Headers {
		fmt.Printf("%s: %s\n", color.CyanString(header.Name), color.YellowString(header.Value))
	}
	if config.ShowBody {
		fmt.Printf("%s:\n", color.CyanString("Body"))
		fmt.Println(string(proxyRequest.Body))
	}
	fmt.Println("")
}

func printResponse(resp *http.Response, body string, showBody bool) {
	color.Blue("Response")
	fmt.Printf("%s: %s\n", color.CyanString("Status"), color.YellowString(resp.Status))
	for name, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", color.CyanString(name), color.YellowString(value))
		}
	}
	if showBody {
		fmt.Printf("%s:\n", color.CyanString("Body"))
		fmt.Println(body)
	}
}
