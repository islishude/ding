package ding

import (
	"testing"
)

func TestValidateToken(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{"64dccf44fac0643d34e461c36175ca109b1a477b208b595268501cc5aef48fe8"}, true},
		{"2", args{"73b868bcdc90e9bf534430cf16503a231299add717c5a02f7a530b49757ec898"}, true},
		{"3", args{"fa53578288687f2ecdb6a10d35e095280343660530bf16a8cfeedca966416c31"}, true},
		{"4", args{"c97fc1833897b44a6cda0b4467ea35f6d7b7769462d67c96faef1e6f5ccc078a"}, true},
		{"5", args{"C97FC1833897B44A6CDA0B4467EA35F6D7B7769462D67C96FAEF1E6F5CCC078A"}, false},
		{"6", args{"ABC123"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateToken(tt.args.key); got != tt.want {
				t.Errorf("ValidateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateSecretKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{"SECd3cf4b5914673e765148426ab2904f795603d6734b458ddb6a53ba5652fb3b27"}, true},
		{"2", args{"SECb89680050aeb698c538322cf11364565ef4553d916d960cc15b0e392db5e2f31"}, true},
		{"3", args{"SECa2588f4e5ce18f1c715a7f9390a2a6d8a9f6f36c6311e5b1165ac81cfb918fdf"}, true},
		{"4", args{"SECa081e0a911f97ff50fbce652a99944215c5db2d95867e65016e43561d74ff188"}, true},
		{"5", args{"C97FC1833897B44A6CDA0B4467EA35F6D7B7769462D67C96FAEF1E6F5CCC078A"}, false},
		{"6", args{"ABC123"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateSecretKey(tt.args.key); got != tt.want {
				t.Errorf("ValidateSecretKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
