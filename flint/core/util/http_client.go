package util

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
)

type HttpError struct {
	Code    int
	Message string
}

func (httpError *HttpError) Error() string {
	return fmt.Sprintf("Code %d: %s", httpError.Code, httpError.Message)
}

type HttpResponse struct {
	StatusCode int
	Body       map[string]interface{}
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

func (httpClient *HttpClient) Request(method string, url string, reader io.Reader, acceptedStatusCodes []int, autohandleErrors bool) (*HttpResponse, error) {
	req, err := http.NewRequest(method, httpClient.BaseUrl+url, reader)

	if err != nil {
		if !autohandleErrors {
			return nil, err
		}
		fmt.Println("failed to create an http request")
		fmt.Println(err)
		os.Exit(-1)
	}

	if httpClient.Headers != nil {
		for k, v := range httpClient.Headers {
			req.Header.Add(k, v)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if !autohandleErrors {
			return nil, err
		}
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if !autohandleErrors {
			return nil, err
		}
		log.Println("Error while reading the response bytes:", err)
	}

	if acceptedStatusCodes == nil {
		acceptedStatusCodes = []int{http.StatusOK, http.StatusCreated}
	}

	if !slices.Contains(acceptedStatusCodes, resp.StatusCode) {
		if !autohandleErrors {
			return nil, &HttpError{
				Code:    resp.StatusCode,
				Message: string(body),
			}
		}
		fmt.Printf("%v request to %v resulted in an unacceptable status code %v (acceptable status codes are %v)\n", method, httpClient.BaseUrl+url, resp.Status, acceptedStatusCodes)
		if resp.StatusCode == 422 {
			panic("try reviewing your manifest")
		}
		var respJson map[string]any
		err := json.Unmarshal(body, &respJson)
		if err != nil {
			fmt.Println(respJson)
			panic(err)
		}
		fmt.Println(respJson["message"].(string) + "\n")
		os.Exit(1)
	}

	var mapBody = make(map[string]interface{})

	err = json.Unmarshal(body, &mapBody)

	if err != nil {
		panic(err)
	}

	return &HttpResponse{
		StatusCode: resp.StatusCode,
		Body:       mapBody,
	}, nil
}

func (httpClient *HttpClient) Post(url string, reader io.Reader, acceptedStatusCodes []int, autohandleErrors bool) (*HttpResponse, error) {
	return httpClient.Request("POST", url, reader, acceptedStatusCodes, autohandleErrors)
}

func (httpClient *HttpClient) Put(url string, reader io.Reader, acceptedStatusCodes []int, autohandleErrors bool) (*HttpResponse, error) {
	return httpClient.Request("PUT", url, reader, acceptedStatusCodes, autohandleErrors)
}

func (httpClient *HttpClient) Delete(url string, acceptedStatusCodes []int, autohandleErrors bool) (*HttpResponse, error) {
	return httpClient.Request("DELETE", url, bytes.NewReader(make([]byte, 0)), acceptedStatusCodes, autohandleErrors)
}

func (httpClient *HttpClient) Get(url string, acceptedStatusCodes []int, autohandleErrors bool) (*HttpResponse, error) {
	return httpClient.Request("GET", url, bytes.NewReader(make([]byte, 0)), acceptedStatusCodes, autohandleErrors)
}
