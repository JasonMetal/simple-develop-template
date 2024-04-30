package number

import "testing"

func TestGenerateSpanId(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			"nornal",
			args{length: 10},
			"1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSpanId(tt.args.length); got != tt.want {
				t.Errorf("GenerateSpanId() = %v, want %v", got, tt.want)
			}
		})
	}
}
