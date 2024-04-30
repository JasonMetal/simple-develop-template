package config

import (
	configLib "github.com/olebedev/config"
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *configLib.Config
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConfig(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetYaml(t *testing.T) {
	type args struct {
		filePath   string
		configData any
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetYaml(tt.args.filePath, tt.args.configData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}
