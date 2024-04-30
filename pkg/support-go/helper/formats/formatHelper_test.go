package formats

import "testing"

func TestCheckMobile(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"纯数字", args{"123abc5623"}, false},
		{"11位", args{"123456"}, false},
		{"第一位要为1", args{"23912345678"}, false},
		{"第二位只能345789", args{"12912345678"}, false},
		{"正确格式", args{"18913245678"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckMobile(tt.args.str); got != tt.want {
				t.Errorf("CheckMobile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckEmail(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "邮箱格式-正确",
			args: args{
				str: "1254@qq.com",
			},
			want: true,
		},
		{
			name: "邮箱格式-错误",
			args: args{
				str: "1254.com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckEmail(tt.args.str); got != tt.want {
				t.Errorf("CheckEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
