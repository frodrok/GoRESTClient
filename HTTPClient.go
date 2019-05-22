package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	//	"strconv"
)

type HttpResponse struct {
	status        int
	statusString  string
	body          string
	contentType   string
	contentLength string
}

type HttpRequest struct {
	method   string
	url      string
	username string
	password string
	body     string
	headers  map[string]string
}

func basicAuth(username string, password string) string {
	var auth = username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type HttpRequestFunction func(*HttpRequest) HttpResponse

func CallHttp(req *HttpRequest) HttpResponse {

	/* Receive a pointer to a HttpRequest - extract data, perform HTTP request
	and return the result as a HttpResponse */
	bodyBuffer := bytes.NewBuffer([]byte(req.body))

	var newRequest, err = http.NewRequest(req.method, req.url, bodyBuffer)

	if req.headers != nil {
		for key, value := range req.headers {
			newRequest.Header.Add(key, value)
		}
	}

	if err != nil {

		return HttpResponse{
			500,
			"err",
			err.Error(),
			"0",
			"0",
		}

	}

	newRequest.Header.Add("Authorization", "Basic "+basicAuth(req.username, req.password))

	resp, err := client.Do(newRequest)

	if err != nil {

		return HttpResponse{
			500,
			"",
			err.Error(),
			"0",
			"0",
		}

	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	//LOG(string(body))
	//LOG(strconv.Itoa(resp.StatusCode))

	return HttpResponse{
		status:        resp.StatusCode,
		statusString:  resp.Status,
		body:          string(body),
		contentType:   resp.Header.Get("Content-Type"),
		contentLength: resp.Header.Get("Content-Length"),
	}
}

func CallHttpSimple(method string, url string) HttpResponse {

	req, err := http.NewRequest(method, url, nil)

	if err != nil {

		return HttpResponse{
			500,
			"",
			"error",
			"0",
			"",
		}

	}

	resp, err := client.Do(req)

	if err != nil {
		return HttpResponse{
			500,
			"",
			"error",
			"0",
			"",
		}
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	for k, v := range resp.Header {
		fmt.Printf("key[%url] value[%url]\n", k, v)
	}

	return HttpResponse{
		status:        resp.StatusCode,
		body:          string(body),
		contentType:   resp.Header.Get("Content-Type"),
		contentLength: resp.Header.Get("Content-Length"),
	}

}
