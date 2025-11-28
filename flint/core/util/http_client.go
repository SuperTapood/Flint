package util

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
)

type HttpResponse struct {
	StatusCode int
	Body       []byte
}

type HttpClient struct {
	Headers map[string]string
	Client  *http.Client
	BaseUrl string
}

func NewHttpClient(headers map[string]string, baseUrl string) *HttpClient {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	return &HttpClient{
		Headers: headers,
		Client:  &http.Client{},
		BaseUrl: baseUrl,
	}
}

func (httpClient *HttpClient) Request(method string, url string, reader io.Reader, acceptedStatusCodes []int) *HttpResponse {
	req, err := http.NewRequest(method, url, reader)

	if err != nil {
		panic(err)
	}

	if httpClient.Headers != nil {
		for k, v := range httpClient.Headers {
			req.Header.Add(k, v)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, er := io.ReadAll(resp.Body)
	if er != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	if acceptedStatusCodes == nil {
		acceptedStatusCodes = []int{http.StatusOK}
	}

	if !slices.Contains(acceptedStatusCodes, resp.StatusCode) {
		fmt.Printf("%v request to %v resulted in an unacceptable status code %v (acceptable status codes are %v)\n\n", method, httpClient.BaseUrl+url, resp.StatusCode, acceptedStatusCodes)
		os.Exit(1)
	}

	return &HttpResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}
}

func (httpClient *HttpClient) Post(url string, reader io.Reader, acceptedStatusCodes []int) *HttpResponse {
	return httpClient.Request("POST", url, reader, acceptedStatusCodes)
}

func (httpClient *HttpClient) Put(url string, reader io.Reader, acceptedStatusCodes []int) *HttpResponse {
	return httpClient.Request("PUT", url, reader, acceptedStatusCodes)
}

func (httpClient *HttpClient) Delete(url string, acceptedStatusCodes []int) *HttpResponse {
	return httpClient.Request("DELETE", url, bytes.NewReader(make([]byte, 0)), acceptedStatusCodes)
}

func (httpClient *HttpClient) Get(url string, acceptedStatusCodes []int) *HttpResponse {
	return httpClient.Request("GET", url, bytes.NewReader(make([]byte, 0)), acceptedStatusCodes)
}
