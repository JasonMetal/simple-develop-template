package strings

import (
	"reflect"
	"testing"
)

func TestBytesToStr(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToStr(tt.args.b); got != tt.want {
				t.Errorf("BytesToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Empty(tt.args.val); got != tt.want {
				t.Errorf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrToBytes(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubStr(t *testing.T) {

	tests := []struct {
		str    string
		offset int
		length int
		want   string
	}{
		// TODO: Add test cases.
		{
			"abcdefg",
			1,
			2,
			"bc",
		},
		{
			"我叫吴伟池",
			1,
			2,
			"叫吴",
		},
		{
			"hello world",
			3,
			3,
			"lo ",
		},
		{
			"hello world",
			0,
			0,
			"",
		},
		{
			"FBL-CXL",
			0,
			7,
			"FBL-CXL",
		},
		{
			"FBL-CXL",
			0,
			-7,
			"",
		},
		{
			"1！*@（……%（*@中文",
			0,
			6,
			"1！*@（…",
		},
	}
	for _, tt := range tests {
		if got := SubStr(tt.str, tt.offset, tt.length); got != tt.want {
			t.Errorf("SubStr() = %v, want %v", got, tt.want)
		}
	}
}
