package main

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"net/http"
)

func printRequest(proxiedUrl string, r *http.Request, showBody bool) {
	fmt.Println("---")
	color.Blue("Request")
	color.Green("%s %s -> %s", r.Method, r.URL.Path, proxiedUrl)
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", color.CyanString(name), color.YellowString(value))
		}
	}
	if showBody {
		fmt.Printf("%s:\n", color.CyanString("Body"))
		body, _ := ioutil.ReadAll(r.Body)
		fmt.Println(string(body))
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
