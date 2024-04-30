package strings

import "testing"

func TestBase64Decode(t *testing.T) {
	type args struct {
		str string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"numbers",
			args{"MTIzNDU="},
			"12345",
		},
		{
			"strnums",
			args{"YWJjZDEyMw=="},
			"abcd123",
		},
		{
			"chinese",
			args{"JXU0RjYwJXU1OTdEJTJDJXU0RTE2JXU3NTRD"},
			"你好,世界",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64Decode(tt.args.str); got != tt.want {
				t.Errorf("Base64Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBase64Encode(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"numbers",
			args{"12345"},
			"MTIzNDU=",
		},
		{
			"strnums",
			args{"abcd123"},
			"YWJjZDEyMw==",
		},
		{
			"chinese",
			args{"你好,世界"},
			"JXU0RjYwJXU1OTdEJTJDJXU0RTE2JXU3NTRD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Base64Encode(tt.args.str); got != tt.want {
				t.Errorf("Base64Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
