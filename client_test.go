package ding

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

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

		if sign := querystring.Get("sign"); sign != "ZkzB968DOpZVkzHPYH0C67nTCmI5V3T41MINQKncc3U%3D" {
			t.Log(sign)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"errcode": 100, "errmsg": "sign error"})
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]interface{}{"errcode": 0})
	})

	handler.HandleFunc("/invalidjson", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("invalid json"))
	})

	server := httptest.NewServer(handler)

	{
		client := &clientimpl{
			api:    server.URL,
			tokens: []AccessToken{{"token", ""}},
			client: http.DefaultClient,
			now:    func() string { return "timetamp" },
		}

		if err := client.request(nil, nil); err == nil { // nolint
			t.Fatalf("want error but got nil error")
		}
	}

	{
		client := &clientimpl{
			api:    server.URL + "/invalidjson",
			tokens: []AccessToken{{"token", ""}},
			client: http.DefaultClient,
			now:    func() string { return "timetamp" },
		}

		if err := client.request(context.Background(), nil); err == nil { // nolint
			t.Fatalf("want error but got nil error")
		}
	}

	{
		client := &clientimpl{
			api:    server.URL,
			tokens: []AccessToken{{"token", "4bb7292e"}},
			client: http.DefaultClient,
			now:    func() string { return "1576759748808" },
		}

		if err := client.request(context.Background(), map[string]interface{}{}); err != nil {
			t.Fatalf("want no error but got error %v", err)
		}
	}

	{
		client := &clientimpl{
			api:    server.URL,
			tokens: []AccessToken{{"token", ""}},
			client: http.DefaultClient,
			now:    func() string { return "timetamp" },
		}

		if err := client.request(context.Background(), nil); err == nil {
			t.Fatalf("want error %s but got error %v", "sign error", err)
		}
	}

	{
		client := &clientimpl{
			api:    server.URL,
			tokens: nil,
			client: http.DefaultClient,
			now:    func() string { return "timetamp" },
		}

		if err := client.request(context.Background(), nil); err == nil {
			t.Fatalf("want error %s but got error %v", "no access token error", err)
		}
	}
}

func Test_New(t *testing.T) {
	var client = New(nil, "token", "hmkey").(*clientimpl)

	if len(client.tokens) != 1 {
		t.Fatal("tokens length should be 1")
	}

	if client.now == nil {
		t.Fatal("now is nil")
	}

	if client.now() == "" {
		t.Fatal("now() isn't correct")
	}

	accessToken := client.tokens[0]
	if accessToken.Token != "token" {
		t.Fatal("endurl isn't correct")
	}
	if accessToken.Key != "hmkey" {
		t.Fatal("hmkey isn't correct")
	}
}

func Test_clientimpl_nextAccessToken(t *testing.T) {
	tests := []struct {
		name   string
		tokens []AccessToken
		want   AccessToken
		want1  bool
	}{
		{"no access tokens", nil, AccessToken{}, false},
		{"only one tokens", []AccessToken{{"1", "1"}}, AccessToken{"1", "1"}, true},
	}
	for _, tt := range tests {
		d := &clientimpl{tokens: tt.tokens}
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := d.nextAccessToken()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("clientimpl.nextAccessToken() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("clientimpl.nextAccessToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_clientimpl_nextAccessToken_multi(t *testing.T) {
	d := &clientimpl{tokens: []AccessToken{{"0", "0"}, {"1", "1"}, {"2", "2"}, {"3", "3"}}}

	{
		i := 0
		current, ok := d.nextAccessToken()
		if !ok {
			t.Fatalf("test %d should true", i)
		}
		want := strconv.Itoa(i)
		if want != current.Key || want != current.Token {
			t.Fatalf("test %d: should %s but got %s/%s", i, want, current.Key, current.Token)
		}

		if d.index != 1 {
			t.Fatalf("test %d: index should be %d but got %d", i, 1, d.index)
		}
	}

	{
		i := 1
		current, ok := d.nextAccessToken()
		if !ok {
			t.Fatalf("test %d should true", i)
		}
		want := strconv.Itoa(i)
		if want != current.Key || want != current.Token {
			t.Fatalf("test %d: should %s but got %s/%s", i, want, current.Key, current.Token)
		}
		if d.index != 2 {
			t.Fatalf("test %d: index should be %d but got %d", i, 2, d.index)
		}
	}

	{
		i := 2
		current, ok := d.nextAccessToken()
		if !ok {
			t.Fatalf("test %d should true", i)
		}
		want := strconv.Itoa(i)
		if want != current.Key || want != current.Token {
			t.Fatalf("test %d: should %s but got %s/%s", i, want, current.Key, current.Token)
		}

		if d.index != 3 {
			t.Fatalf("test %d: index should be %d but got %d", i, 3, d.index)
		}
	}

	{
		i := 3
		current, ok := d.nextAccessToken()
		if !ok {
			t.Fatalf("test %d should true", i)
		}
		want := strconv.Itoa(i)
		if want != current.Key || want != current.Token {
			t.Fatalf("index %d: should %s but got %s/%s", i, want, current.Key, current.Token)
		}

		if d.index != 4 {
			t.Fatalf("test %d: index should be %d but got %d", i, 4, d.index)
		}
	}

	{
		i := 4
		current, ok := d.nextAccessToken()
		if !ok {
			t.Fatalf("test %d should true", i)
		}
		want := "0"
		if want != current.Key || want != current.Token {
			t.Fatalf("test %d: should %s but got %s/%s", i, want, current.Key, current.Token)
		}

		if d.index != 1 {
			t.Fatalf("test %d: index should be %d but got %d", i, 1, d.index)
		}
	}

	{
		i := 5
		current, ok := d.nextAccessToken()
		if !ok {
			t.Fatalf("test %d: should true", i)
		}
		want := "1"
		if want != current.Key || want != current.Token {
			t.Fatalf("test %d: should %s but got %s/%s", i, want, current.Key, current.Token)
		}

		if d.index != 2 {
			t.Fatalf("test %d: index should be %d but got %d", i, 2, d.index)
		}
	}
}
