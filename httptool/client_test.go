package httptool

import (
	"testing"
)

type ClientTest struct {
	Id        string         `json:"id" in:"path"`
	Type      string         `json:"type" in:"query"`
	Param     ClientTestBody `json:"Param" in:"body"`
	Debug     string         `json:"debug" in:"cookie"`
	RequestId string         `json:"X-Request-ID" in:"header"`
}

type ClientTestBody struct {
	ExOrderNo string `json:"exOrderNo"`
	Amount    int64  `json:"Amount"`
}

type ClientTestBankListVO struct {
	Code string `json:"code"` //
	Name string `json:"name"` //
}

func TestClient_Request(t *testing.T) {
	client := Client{
		Host: "192.168.1.116:8080",
		Port: 8080,
	}

	req := ClientTest{
		Id:   "123",
		Type: "bill",
		Param: ClientTestBody{
			ExOrderNo: "ef7d78a5-7af3-4027-bacc-6e4b95f2b283",
			Amount:    1000,
		},
		Debug:     "abcxxx-77e1c83b-7bb0-437b-bc50-a7a58e5660ac",
		RequestId: "77e1c83b-7bb0-437b-bc50-a7a58e5660ac",
	}
	var body []byte
	var err error
	body, err = client.Request("GET", "/fpx/api/v1.0/bank-list", req)
	if err != nil {
		t.Fatal(err)
	}
	resp := []ClientTestBankListVO{}
	err = client.ParseToResult(body, &resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
