package configuration_test

import (
	"testing"
)

func TestLoad(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "0. empty path",
			args: args{
				path: "",
			},
			wantErr: true,
		},
		{
			name: "1. valid path",
			args: args{
				path: "./../../../test/source/config-valid.yaml",
			},
			wantErr: false,
		},
		{
			name: "2. invalid path",
			args: args{
				path: "./../../../test/source/nofile.yaml",
			},
			wantErr: true,
		},
		{
			name: "3. invalid yaml config values",
			args: args{
				path: "./../../../test/source/config-invalid.yaml",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			configuration.Reset()
			if err := configuration.Load(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoad_Reload(t *testing.T) {
	configuration.Reset()

	loadingErr := configuration.Load("./../../../test/source/config-valid.yaml")
	if loadingErr != nil {
		t.Errorf("configuration.Load() err = %+v, wantErr = false", loadingErr)
	}

	loadingErr = configuration.Load("./../../../test/source/config-valid.yaml")
	if loadingErr == nil {
		t.Errorf("recall configuration.Load() err = %+v, wantErr = true", loadingErr)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		loadingErr    bool
		validatingErr bool
	}{
		{
			name:          "0. invalid yaml configuration values",
			path:          "./../../../test/source/config-invalid.yaml",
			loadingErr:    true,
			validatingErr: true,
		},
		{
			name:          "0. valid yaml configuration values",
			path:          "./../../../test/source/config-valid.yaml",
			loadingErr:    false,
			validatingErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			configuration.Reset()
			loadingErr := configuration.Load(tt.path)
			if (loadingErr != nil) != tt.loadingErr {
				t.Errorf("\nLoad() error = \n%v \n\nwantErr = \n%v", tt.loadingErr, loadingErr)
			}
			validatingErr := configuration.Validate()
			if (validatingErr != nil) != tt.validatingErr {
				t.Errorf("\nValidate() error = \n%v \n\nwantErr = \n%v", tt.validatingErr, validatingErr)
			}
		})
	}
}
