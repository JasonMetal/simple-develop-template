package maps

import (
	"reflect"
	"testing"
)

func TestClear(t *testing.T) {
	type args struct {
		m M
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Clear(tt.args.m)
		})
	}
}

func TestClone(t *testing.T) {
	type args struct {
		m M
	}
	tests := []struct {
		name string
		args args
		want M
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clone(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	type args struct {
		dst M
		src M
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Copy(tt.args.dst, tt.args.src)
		})
	}
}

func TestDeleteFunc(t *testing.T) {
	type args struct {
		m   M
		del func(K, V) bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DeleteFunc(tt.args.m, tt.args.del)
		})
	}
}

func TestEqual(t *testing.T) {
	type args struct {
		m1 M1
		m2 M2
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
			if got := Equal(tt.args.m1, tt.args.m2); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualFunc(t *testing.T) {
	type args struct {
		m1 M1
		m2 M2
		eq func(V1, V2) bool
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
			if got := EqualFunc(tt.args.m1, tt.args.m2, tt.args.eq); got != tt.want {
				t.Errorf("EqualFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	type args struct {
		m M
	}
	tests := []struct {
		name string
		args args
		want []K
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Keys(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValues(t *testing.T) {
	type args struct {
		m M
	}
	tests := []struct {
		name string
		args args
		want []V
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Values(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}
