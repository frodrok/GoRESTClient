package httpClient

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	//	"strconv"
)

var client = &http.Client{}

type HttpResponse struct {
	Status        int
	StatusString  string
	Body          string
	ContentType   string
	ContentLength string
}

type HttpRequest struct {
	Method   string
	Url      string
	Username string
	Password string
	Body     string
	Headers  map[string]string
}

func basicAuth(username string, password string) string {
	var auth = username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

type HttpRequestFunction func(*HttpRequest) HttpResponse

func CallHttp(req *HttpRequest) HttpResponse {

	/* Receive a pointer to a HttpRequest - extract data, perform HTTP request
	and return the result as a HttpResponse */
	bodyBuffer := bytes.NewBuffer([]byte(req.Body))

	var newRequest, err = http.NewRequest(req.Method, req.Url, bodyBuffer)

	if req.Headers != nil {
		for key, value := range req.Headers {
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

	newRequest.Header.Add("Authorization", "Basic "+basicAuth(req.Username, req.Password))

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
		Status:        resp.StatusCode,
		StatusString:  resp.Status,
		Body:          string(body),
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.Header.Get("Content-Length"),
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
		Status:        resp.StatusCode,
		Body:          string(body),
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.Header.Get("Content-Length"),
	}

}
