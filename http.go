package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type ProxyHeader struct {
	Name  string
	Value string
}

type ProxyRequest struct {
	Method      string
	Path        string
	Headers     []ProxyHeader
	Body        []byte
	ProxiedHost string
}

type ProxyResponse struct {
	Status  string
	Headers []ProxyHeader
	Body    []byte
}

func NewProxyRequestFrom(httpRequest *http.Request, proxiedHost string) ProxyRequest {
	body, _ := ioutil.ReadAll(httpRequest.Body)
	headers := make([]ProxyHeader, 0)
	for name, values := range httpRequest.Header {
		for _, value := range values {
			headers = append(headers, ProxyHeader{name, value})
		}
	}
	return ProxyRequest{
		Method:      httpRequest.Method,
		Path:        httpRequest.URL.Path,
		Body:        body,
		Headers:     headers,
		ProxiedHost: proxiedHost,
	}
}

func (request *ProxyRequest) ProxiedUrl() string {
	return request.ProxiedHost + request.Path
}

func (request *ProxyRequest) NewProxiedHttpRequest() *http.Request {
	httpRequest, _ := http.NewRequest(request.Method, request.ProxiedUrl(), bytes.NewReader(request.Body))
	return httpRequest
}

func NewProxyResponseFrom(httpResponse *http.Response) ProxyResponse {
	body, _ := ioutil.ReadAll(httpResponse.Body)
	headers := make([]ProxyHeader, 0)
	for name, values := range httpResponse.Header {
		for _, value := range values {
			headers = append(headers, ProxyHeader{name, value})
		}
	}
	return ProxyResponse{
		Status:  httpResponse.Status,
		Body:    body,
		Headers: headers,
	}
}
