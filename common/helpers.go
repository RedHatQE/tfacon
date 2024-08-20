package common

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// PrintGreen is a helper function that
// prints str in green to terminal.
func PrintGreen(str string) {
	color.Green(str)
}

// StructToString converts any struct to a string
func StructToString(i interface{}) string {
	var result string
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name
		result += fmt.Sprintf("%s: \t %v\n", fieldName, field.Interface())
	}
	return result
}

// FileExist is a helper function that
// tell if a file exites.
func FileExist(str string) bool {
	if _, err := os.Stat(str); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

// PrintRed is a helper function that
// prints str in red to terminal.
func PrintRed(str string) {
	color.Red(str)
}

// PrintHeader is a helper function
// for the whole program to print
// header information.
func PrintHeader(version string) {
	fmt.Println("--------------------------------------------------")
	fmt.Printf("tfacon  %s\n", version)
	fmt.Println("Copyright (C) 2023, Red Hat, Inc.")
	fmt.Print("-------------------------------------------------\n\n\n")
}

// SendHTTPRequest is a helper function that
// deals with all http operation for tfacon.
func SendHTTPRequest(ctx context.Context, method, url,
	authToken string, body io.Reader, client *http.Client) ([]byte, bool, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)

	if authToken != "" {
		req.Header.Add("Authorization", fmt.Sprintf("bearer %s", authToken))
	}

	if err != nil {
		err = fmt.Errorf("tfacon http handler crashed, request built failed, could be a bad request: %w", err)

		return nil, false, err
	}

	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)

	defer func() {
		if req != nil && req.Body != nil {
			err = req.Body.Close()
			HandleError(err, "nopanic")
		}

		if resp != nil && resp.Body != nil {
			err = resp.Body.Close()
			HandleError(err, "nopanic")
		}
	}()

	if err != nil {
		err = fmt.Errorf("tfacon http handler crashed:%w", err)

		return nil, false, err
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("tfacon http handler crashed, response read failed:%w", err)

		return nil, false, err
	}

	success, err := httpHelper(method, resp)

	return d, success, err
}

func httpHelper(method string, resp *http.Response) (bool, error) {
	var err error
	if resp.StatusCode == http.StatusOK {
		return true, err
	}

	if method == "POST" && resp.StatusCode == http.StatusCreated {
		return true, err
	}

	err = fmt.Errorf("http handler request exception, status code is:%d, err is %w", resp.StatusCode, err)

	return false, err
}

// HandleError is the Error handler
// for the whole tfacon.
func HandleError(err error, method string) {
	if err != nil {
		// print out the caller information
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			fmt.Printf("called from %s\n", details.Name())
		}
		if method == "panic" {
			panic(err)

		} else {
			fmt.Println(err)
		}
	}
}

// InitTFAConfigFile read the TFA config file for config
func InitTFAConfigFile(viper *viper.Viper) {
	var file []byte

	var err error
	exist := FileExist("./tfacon.cfg")
	if os.Getenv("TFACON_CONFIG_PATH") != "" {
		file, err = ioutil.ReadFile(os.Getenv("TFACON_CONFIG_PATH"))
	} else if exist {
		file, err = ioutil.ReadFile("./tfacon.cfg")
	} else {
		_, err = os.Create("tfacon.cfg")
		HandleError(err, "nopanic")
		file, err = ioutil.ReadFile("./tfacon.cfg")
		HandleError(err, "nopanic")
	}
	HandleError(err, "nopanic")
	viper.SetConfigType("ini")
	viper.SetDefault("config.concurrency", true)
	viper.SetDefault("config.auto_finalize_defect_type", false)
	viper.SetDefault("config.auto_finalization_threshold", 0.5)
	viper.SetDefault("config.retry_times", 20)
	viper.SetDefault("config.add_attributes", false)
	err = viper.ReadConfig(bytes.NewBuffer(file))
	HandleError(err, "nopanic")
}

// UpdateTFAConfig update the TFA config and return new TFAConfig
func UpdateTFAConfig(concurrent bool, addAttributes bool, enableRE bool, autoFinalizeDefectType bool, autoFinalizationThreshold float32, retryTimes int, verbose bool) *TFAConfig {
	tfaConfig := &TFAConfig{
		Concurrent:                concurrent,
		AddAttributes:             addAttributes,
		Re:                        enableRE,
		AutoFinalizeDefectType:    autoFinalizeDefectType,
		AutoFinalizationThreshold: autoFinalizationThreshold,
		RetryTimes:                retryTimes,
		Verbose:                   verbose,
	}
	return tfaConfig
}
