package di

import (
	"github.com/chandy20/prueba-smartjobandina/beer/find/v1/internal/ctx"
	"github.com/chandy20/prueba-smartjobandina/beer/repository"
	"os"
	"reflect"
	"testing"
)

func Test_providerAwsSession(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should_build_aws_config_correctly",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := providerAWSConfig()
			if tt.wantErr && err == nil {
				t.Errorf("awsConfigProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_providerBeerRepository(t *testing.T) {
	tests := []struct {
		name       string
		want       *repository.BeerRepository
		setEnvVars func()
		wantErr    bool
	}{
		{
			name:       "should_fail_because_env_var_is_not_defined",
			want:       nil,
			setEnvVars: func() {},
			wantErr:    true,
		},
		{
			name: "should_build_repository_correctly",
			want: repository.NewBeerRepository(nil, "some-table", nil),
			setEnvVars: func() {
				err := os.Setenv("DYNAMODB_BEERS", "some-table")
				if err != nil {
					t.Errorf("error setting env var %v", err)
				}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setEnvVars()
			got, err := providerBeerRepository(nil, nil)
			if !tt.wantErr && err != nil {
				t.Errorf("handlerProvider() got = %v, want %v", got, tt.want)
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("providerBeerRepository() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func Test_provideNewHandler(t *testing.T) {
	tests := []struct {
		name string
		want *ctx.Handler
	}{
		{
			name: "should_build_handler_successfully",
			want: ctx.NewHandler(nil, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := provideNewHandler(nil, nil)

			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("provideNewHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
