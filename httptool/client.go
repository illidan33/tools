package httptool

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"reflect"
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

func (client *Client) Request(method string, uri string, req interface{}, res interface{}, headers ...map[string]string) (err error) {
	client.initConfig()
	if len(headers) == 0 {
		headers = []map[string]string{
			{
				fasthttp.HeaderContentType: "application/json;charset=utf-8",
			},
		}
	}
	request := HttpRequest{
		Url:     fmt.Sprintf("%s://%s/%s", client.Mode, client.Host, uri),
		Method:  method,
		Timeout: client.TimeOut,
		Params:  req,
	}
	if len(headers) > 0 {
		for _, header := range headers {
			request.AddHeaders(header)
		}
	}
	var body []byte
	body, err = request.Do()
	if err != nil {
		return
	}

	rs := Result{}
	err = json.Unmarshal(body, &rs)
	if err != nil {
		return err
	}
	if rs.Code != 0 {
		return errors.New(fmt.Sprintf("Code: %d, Msg: %s, Time: %s", rs.Code, rs.Msg, time.Unix(rs.Time, 0).Format("2006-01-02 15:04:05")))
	}

	if res != nil {
		rf := reflect.TypeOf(res)
		if rf.Kind() == reflect.Struct || rf.Kind() == reflect.Slice || rf.Kind() == reflect.Array {
			err = json.Unmarshal(rs.Data, res)
			if err != nil {
				return
			}
		} else {
			res = rs.Data
		}
	}

	return
}
