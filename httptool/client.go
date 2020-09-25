package httptool

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"myprojects/tools/common"
	"myprojects/tools/gen/util/types"
	"reflect"
	"strings"
	"time"
)

type Client struct {
	// for config
	Host    string
	Port    int32
	Mode    string
	TimeOut time.Duration
}

func (client *Client) initConfig() *Client {
	if client.Host == "" {
		client.Host = "127.0.0.1"
	}
	if client.Port == 0 {
		client.Port = 80
	}
	if client.Mode == "" {
		client.Mode = "http"
	}
	if client.TimeOut == 0 {
		client.TimeOut = time.Second * 5
	}
	return client
}

func (client *Client) parsePathParams(uri string, params map[string]string) (url string) {
	url = fmt.Sprintf("%s://%s/%s", client.Mode, client.Host, uri)
	if len(params) > 0 {
		for s, pm := range params {
			url = strings.ReplaceAll(url, "{"+s+"}", pm)
		}
	}
	return
}

func (client *Client) parseQueryParams(url string, params map[string]string) string {
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
	return strings.Trim(urlAddress, "&")
}

func (client *Client) parseHeaderParams(headerParams map[string]string, cookieParams map[string]string) (rs map[string]string) {
	rs = headerParams
	if len(cookieParams) > 0 {
		s := ""
		for n, v := range cookieParams {
			if s == "" {
				s = n + "=" + v
			} else {
				s = fmt.Sprintf("%s; %s=%s", s, n, v)
			}
		}
		rs["Cookie"] = s
	}

	return
}

func (client *Client) parseParams(uri string, params interface{}) (newUrl string, req interface{}, headers map[string]string, err error) {
	if params == nil {
		return
	}
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	headerParams := map[string]string{}
	cookieParams := map[string]string{}

	pmType := reflect.TypeOf(params)
	if pmType.Kind() != reflect.Struct {
		err = errors.New("params need struct")
		return
	}
	pmValue := reflect.ValueOf(params)

	n := pmType.NumField()
	for i := 0; i < n; i++ {
		f := pmType.Field(i)
		tag := f.Tag
		in := tag.Get("in")
		if in == "" {
			err = errors.New("field has no tag 'in': " + f.Name)
			return
		}
		name := tag.Get("json")
		if name == "" {
			name = common.ToLowerCamelCase(f.Name)
		}
		switch in {
		case types.SWAGGER_TYPE__BODY:
			req = pmValue.Field(i).Interface()
		case types.SWAGGER_TYPE__PATH:
			pathParams[name] = pmValue.Field(i).String()
		case types.SWAGGER_TYPE__QUERY:
			queryParams[name] = pmValue.Field(i).String()
		case types.SWAGGER_TYPE__HEADER:
			headerParams[name] = pmValue.Field(i).String()
		case types.SWAGGER_TYPE__COOKIE:
			cookieParams[name] = pmValue.Field(i).String()
		default:
			err = errors.New("field tag has wrong 'in': " + in)
		}
	}

	newUrl = client.parsePathParams(uri, pathParams)
	newUrl = client.parseQueryParams(newUrl, queryParams)
	headers = client.parseHeaderParams(headerParams, cookieParams)

	return
}

func (client *Client) Request(method string, uri string, req interface{}, headers ...map[string]string) (result []byte, err error) {
	client.initConfig()

	var url string
	headerMap := map[string]string{}
	var newReq interface{}
	if req != nil {
		url, newReq, headerMap, err = client.parseParams(uri, req)
		if err != nil {
			return
		}
	} else {
		url = fmt.Sprintf("%s://%s/%s", client.Mode, client.Host, uri)
	}
	if _, ok := headerMap[fasthttp.HeaderContentType]; !ok {
		headerMap[fasthttp.HeaderContentType] = "application/json;charset=utf-8"
	}

	request := NewHttpRequest(url, newReq).SetMethod(method).SetTimeout(client.TimeOut)
	if len(headers) > 0 {
		for _, header := range headers {
			request.SetHeaders(header)
		}
	}
	if len(headerMap) > 0 {
		request.SetHeaders(headerMap)
	}

	result, err = request.Do()
	if err != nil {
		return
	}

	return
}

func (client *Client) ParseToResult(body []byte, res interface{}) (err error) {
	rs := Result{}
	err = json.Unmarshal(body, &rs)
	if err != nil {
		return
	}
	if rs.Code != RESULT_SUCCESS {
		return errors.New(fmt.Sprintf("Code: %d, Msg: %s, Time: %s", rs.Code, rs.Msg, time.Unix(rs.Time, 0).Format("2006-01-02 15:04:05")))
	}

	if res != nil {
		rf := reflect.TypeOf(res)
		if rf.Kind() == reflect.Ptr {
			rf = rf.Elem()
		}
		if rf.Kind() == reflect.Struct || rf.Kind() == reflect.Slice || rf.Kind() == reflect.Array {
			err = json.Unmarshal(rs.Data, res)
			if err != nil {
				return
			}
		} else {
			res = rs.Data
		}
	}

	return nil
}
