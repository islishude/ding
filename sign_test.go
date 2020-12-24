package ding

import "testing"

func TestSign(t *testing.T) {
	type fields struct {
		key string
		now string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"1", fields{"036fd40e", "1576758870133"}, "%2BlymVhtFUuawBqgs2WJjcprk3K8HkX2gJkMrXE4eW74%3D"},
		{"2", fields{"2146b67d", "1576758871138"}, "N80sNHfw60EN03Mrl5W7tWpl63N0GTMSHw0%2FS9VtGVc%3D"},
		{"3", fields{"4f7768e1", "1576758872140"}, "9EmKMljsF%2BavVU9624KVA50%2FazwCYeW%2B%2FvJ0puBIvuM%3D"},
		{"4", fields{"d4874a3e", "1576758873145"}, "1Qts2DgLlYRElyavx8m5AsDLbcZ04VLfFQMEvDoFuUU%3D"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getsign(tt.fields.key, tt.fields.now); got != tt.want {
				t.Errorf("getSign() = %v, want %v", got, tt.want)
			}
		})
	}
}
