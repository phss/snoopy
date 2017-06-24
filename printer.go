package main

import (
	"fmt"
	"github.com/fatih/color"
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

func printResponse(response ProxyResponse, config Config) {
	color.Blue("Response")
	fmt.Printf("%s: %s\n", color.CyanString("Status"), color.YellowString(response.Status))
	for _, header := range response.Headers {
		fmt.Printf("%s: %s\n", color.CyanString(header.Name), color.YellowString(header.Value))
	}
	if config.ShowBody {
		fmt.Printf("%s:\n", color.CyanString("Body"))
		fmt.Println(string(response.Body))
	}
}
