package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/RedHatQE/tfacon/common"
	"github.com/RedHatQE/tfacon/connectors"
	"github.com/spf13/viper"
)

func buildTestingViper(configName, configPath string) *viper.Viper {
	var vipertfaconfortest *viper.Viper = viper.New()
	vipertfaconfortest.SetConfigName(configName)
	vipertfaconfortest.AddConfigPath(configPath)
	err := vipertfaconfortest.ReadInConfig()
	common.HandleError(err)
	rpcon := buildEmptyRPCon()
	err = vipertfaconfortest.Unmarshal(rpcon)
	common.HandleError(err)
	return vipertfaconfortest
}
func buildEmptyRPCon() *connectors.RPConnector {
	rpcon := &connectors.RPConnector{Client: &http.Client{}}
	return rpcon
}
func testDataPath() string {
	return "../test_data/workspace_data/"
}
func rpconWithAllParameters() *connectors.RPConnector {
	rpcon := buildEmptyRPCon()
	viper := buildTestingViper("tfacon", testDataPath())
	err := viper.Unmarshal(rpcon)
	common.HandleError(err)
	return rpcon
}
func TestGetCon(t *testing.T) {
	type args struct {
		viper *viper.Viper
	}
	tests := []struct {
		name string
		args args
		want TFACon
	}{
		// TODO: Add test cases.
		{
			name: "test tfacon get with rpcon",
			args: args{
				viper: buildTestingViper("tfacon", testDataPath()),
			},
			want: rpconWithAllParameters(),
		},
		{
			name: "test tfacon get default con",
			args: args{
				viper: buildTestingViper("tfacon_default_rpcon", testDataPath()),
			},
			want: rpconWithAllParameters(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCon(tt.args.viper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCon() = doesn't match the wanted results, please check")
			} else {
				t.Logf("Test %s passed!", tt.name)
			}
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		viperRun    *viper.Viper
		viperConfig *viper.Viper
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.viperRun, tt.args.viperConfig)
		})
	}
}
func viperConfig() *viper.Viper {

	var file []byte
	var err error

	file, err = ioutil.ReadFile(testDataPath() + "tfacon.cfg")

	defer func() {
		if r := recover(); r != nil {
			_, err = os.Create("tfacon.cfg")
			common.HandleError(err)
		}
	}()
	common.HandleError(err)
	viperConfig := viper.New()
	viperConfig.SetConfigType("ini")
	viperConfig.SetDefault("config.concurrency", true)
	viperConfig.SetDefault("config.retry_times", 1)
	viperConfig.SetDefault("config.add_attributes", false)
	err = viperConfig.ReadConfig(bytes.NewBuffer(file))
	common.HandleError(err)
	return viperConfig
}

func TestGetInfo(t *testing.T) {
	type args struct {
		viper *viper.Viper
	}
	tests := []struct {
		name string
		args args
		want TFACon
	}{
		// TODO: Add test cases.		// TODO: Add test cases.
		{
			name: "test tfacon get info with rpcon",
			args: args{
				viper: buildTestingViper("tfacon", testDataPath()),
			},
			want: rpconWithAllParameters(),
		},
		{
			name: "test tfacon get info default con",
			args: args{
				viper: buildTestingViper("tfacon_default_rpcon", testDataPath()),
			},
			want: rpconWithAllParameters(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetInfo(tt.args.viper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInfo() doesn't match the wanted results, please check")
			}
		})
	}
}

func TestValidate(t *testing.T) {
	viperForTestValidate := viperConfig()
	conInvalid := GetCon(buildTestingViper("tfacon_validate", testDataPath()))
	conValid := GetCon(buildTestingViper("tfacon_validate", testDataPath()))
	type args struct {
		con   TFACon
		viper *viper.Viper
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr error
	}{
		{
			name: "Test validate with validate tfacon.yml",
			args: args{
				con:   conInvalid,
				viper: viperForTestValidate,
			},
			want:    false,
			wantErr: fmt.Errorf("validate error: "),
		},
		{
			name: "Test validate with validate with real data, it should pass the validation",
			args: args{
				con:   conValid,
				viper: viperForTestValidate,
			},
			want:    false,
			wantErr: fmt.Errorf("validate error: "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Validate(tt.args.con, tt.args.viper)
			if err != nil {
				if !strings.Contains(err.Error(), tt.wantErr.Error()) {

					t.Errorf("Validate() error = %v\n, wantErr %v\n", err, tt.wantErr)
				}
				return
			}
			if got != tt.want {
				t.Errorf("Validate() doesn't match the wanted results, please check")
			}
		})
	}
}

func Test_runHelper(t *testing.T) {
	type args struct {
		viperConfig *viper.Viper
		ids         []string
		con         TFACon
		operation   string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runHelper(tt.args.viperConfig, tt.args.ids, tt.args.con, tt.args.operation)
		})
	}
}
