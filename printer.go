package main

import (
	"fmt"
	"github.com/fatih/color"
)

func printRequest(request ProxyRequest, config Config) {
	fmt.Println("---")
	color.Blue("Request")
	color.Green("%s %s -> %s", request.Method, request.Path, request.ProxiedUrl())
	printHeaders(request.Headers)
	printBody(request.Body, config)
	fmt.Println("")
}

func printResponse(response ProxyResponse, config Config) {
	color.Blue("Response")
	fmt.Printf("%s: %s\n", color.CyanString("Status"), color.YellowString(response.Status))
	printHeaders(response.Headers)
	printBody(response.Body, config)
}

func printHeaders(headers []ProxyHeader) {
	for _, header := range headers {
		fmt.Printf("%s: %s\n", color.CyanString(header.Name), color.YellowString(header.Value))
	}
}

func printBody(body []byte, config Config) {
	if config.ShowBody {
		fmt.Printf("%s:\n", color.CyanString("Body"))
		fmt.Println(string(body))
	}
}
