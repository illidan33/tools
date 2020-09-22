package http

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
	"github.com/valyala/fasthttp"
)

type TransportWrapper func(rt http.RoundTripper) http.RoundTripper

type HttpRequest struct {
	BaseURL  string
	Method   string
	URI      string
	Headers Metadata
	Timeout  time.Duration
	Req      interface{}
}

type BaseResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
	Time int64           `json:"time"`
}

func (httpRequest *HttpRequest) Do() (result BaseResponse,err error) {
	result = BaseResponse{}
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	for s, v := range httpRequest.Headers {
		for _, v2 := range v {
			req.Header.Set(s, v2)
		}
	}
	req.Header.SetMethod(httpRequest.Method)
	urlAddress := fmt.Sprintf("%s/%s", httpRequest.BaseURL, httpRequest.URI)
	resp := fasthttp.AcquireResponse()
	req.SetRequestURI(urlAddress)
	defer fasthttp.ReleaseResponse(resp)
	if err := fasthttp.DoTimeout(req, resp, httpRequest.Timeout); err != nil {
		fmt.Printf("Http Request Do Error %s", err.Error())
		return
	}
	respBody := resp.Body()
	err = json.Unmarshal(respBody, &result)

	return
}
