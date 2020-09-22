package httptool

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/valyala/fasthttp"
)

type HttpRequest struct {
	Url     string
	Method  string
	Headers map[string]string
	Timeout time.Duration
	Params  interface{}
}

func InitHttpRequest(url string) *HttpRequest {
	return &HttpRequest{
		Url: url,
	}
}

func (httpRequest *HttpRequest) AddHeaders(headers map[string]string) {
	for s, v := range headers {
		httpRequest.Headers[s] = v
	}
}

func (httpRequest *HttpRequest) Do() (result []byte, err error) {
	if httpRequest.Url == "" {
		return nil, errors.New("url should not be empty")
	}
	if httpRequest.Method == "" {
		httpRequest.Method = fasthttp.MethodGet
	}
	if httpRequest.Timeout == 0 {
		httpRequest.Timeout = time.Second * 5
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	for s, v := range httpRequest.Headers {
		req.Header.Set(s, v)
	}
	req.Header.SetMethod(httpRequest.Method)
	url := httpRequest.Url
	if httpRequest.Method == fasthttp.MethodGet && httpRequest.Params != nil {
		url, err = httpRequest.getRequestURL()
		if err != nil {
			return
		}
	} else if httpRequest.Params != nil {
		body, err := json.Marshal(httpRequest.Params)
		if err != nil {
			return nil, err
		}
		req.SetBody(body)
	}
	req.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if e := fasthttp.DoTimeout(req, resp, httpRequest.Timeout); e != nil {
		err = e
		return
	}
	result = resp.Body()

	return
}

func (httpRequest *HttpRequest) Get() (result []byte, err error) {
	httpRequest.Method = fasthttp.MethodGet
	return httpRequest.Do()
}

func (httpRequest *HttpRequest) Post() (result []byte, err error) {
	httpRequest.Method = fasthttp.MethodPost
	return httpRequest.Do()
}

//append request url
func (httpRequest *HttpRequest) getRequestURL() (url string, err error) {
	params, ok := httpRequest.Params.(map[string]string)
	if !ok {
		return "", errors.New("request param should be map")
	}

	var urlAddress = ""
	lastCharctor := url[len(url)-1:]
	if lastCharctor == "?" {
		urlAddress = url + urlAddress
	} else {
		urlAddress = url + "?" + urlAddress
	}
	for k, v := range params {
		if len(k) != 0 && len(v) != 0 {
			urlAddress = urlAddress + k + "=" + v + "&"
		}
	}
	return urlAddress, nil
}
