package ding

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_clientimpl_getSign(t *testing.T) {
	type fields struct {
		hmkey string
		now   func() string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"1", fields{"036fd40e", func() string { return "1576758870133" }}, "&timestamp=1576758870133&sign=%2BlymVhtFUuawBqgs2WJjcprk3K8HkX2gJkMrXE4eW74%3D"},
		{"2", fields{"2146b67d", func() string { return "1576758871138" }}, "&timestamp=1576758871138&sign=N80sNHfw60EN03Mrl5W7tWpl63N0GTMSHw0%2FS9VtGVc%3D"},
		{"3", fields{"4f7768e1", func() string { return "1576758872140" }}, "&timestamp=1576758872140&sign=9EmKMljsF%2BavVU9624KVA50%2FazwCYeW%2B%2FvJ0puBIvuM%3D"},
		{"4", fields{"d4874a3e", func() string { return "1576758873145" }}, "&timestamp=1576758873145&sign=1Qts2DgLlYRElyavx8m5AsDLbcZ04VLfFQMEvDoFuUU%3D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &clientimpl{
				hmkey:  tt.fields.hmkey,
				bhmkey: []byte(tt.fields.hmkey),
				now:    tt.fields.now,
			}
			if got := d.getSign(); got != tt.want {
				t.Errorf("clientimpl.getSign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clientimpl_request(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"errcode": 100, "errmsg": "NOT POST"})
			return
		}

		if req.Header.Get("Content-Type") != "application/json" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"errcode": 100, "errmsg": "content-type error"})
			return
		}

		querystring := req.URL.Query()
		if querystring.Get("timestamp") != "1576759748808" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"errcode": 100, "errmsg": "timestamp error"})
			return
		}

		if querystring.Get("access_token") != "token" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"errcode": 100, "errmsg": "access_token error"})
			return
		}

		if sign := querystring.Get("sign"); sign != "ZkzB968DOpZVkzHPYH0C67nTCmI5V3T41MINQKncc3U=" {
			println(sign)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"errcode": 100, "errmsg": "sign error"})
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]interface{}{"errcode": 0})
	})
	server := httptest.NewServer(handler)

	{
		client := &clientimpl{
			endurl: server.URL + "/?access_token=token",
			hmkey:  "4bb7292e",
			bhmkey: []byte("4bb7292e"),
			client: http.DefaultClient,
			now:    func() string { return "1576759748808" },
		}

		if err := client.request(context.Background(), map[string]interface{}{}); err != nil {
			t.Fatalf("want no error but got error %v", err)
		}
	}

	{
		client := &clientimpl{
			endurl: server.URL + "/?access_token=token",
			hmkey:  "4bb7292e",
			bhmkey: []byte("4bb7292e"),
			client: http.DefaultClient,
			now:    func() string { return "timetamp" },
		}

		if err := client.request(context.Background(), nil); err == nil {
			t.Fatalf("want error %s but got error %v", "sign error", err)
		}
	}

}

func Test_New(t *testing.T) {
	var client = New(nil, "token", "hmkey").(*clientimpl)

	if client.endurl != "https://oapi.dingtalk.com/robot/send?access_token=token" {
		t.Fatal("endurl isn't correct")
	}

	if !reflect.DeepEqual(client.client, http.DefaultClient) {
		t.Fatal("client isn't correct")
	}

	if client.hmkey != "hmkey" {
		t.Fatal("hmkey isn't correct")
	}

	if bytes.Equal([]byte("bhmkey"), client.bhmkey) {
		t.Fatal("bhmkey isn't correct")
	}

	if client.now == nil {
		t.Fatal("now is nil")
	}

	if client.now() == "" {
		t.Fatal("now() isn't correct")
	}
}
