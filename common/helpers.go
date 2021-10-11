package common

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

// PrintGreen is a helper function that
// prints str in green to terminal.
func PrintGreen(str string) {
	color.Green(str)
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
	fmt.Println("Copyright (C) 2021, Red Hat, Inc.")
	fmt.Print("-------------------------------------------------\n\n\n")
}

// SendHTTPRequest is a helper function that
// deals with all http operation for tfacon.
func SendHTTPRequest(ctx context.Context, method, url,
	auth_token string, body io.Reader, client *http.Client) ([]byte, bool, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", auth_token))

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
			HandleError(err)
		}

		if resp != nil && resp.Body != nil {
			err = resp.Body.Close()
			HandleError(err)
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

	// if resp.StatusCode == http.StatusOK {
	// 	return d, true, err
	// }

	// if method == "POST" && resp.StatusCode == http.StatusCreated {
	// 	return d, true, err
	// }

	// err = fmt.Errorf("http handler request exception, status code is:%d, err is %w\n", resp.StatusCode, err)
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

	err = fmt.Errorf("http handler request exception, status code is:%d, err is %w\n", resp.StatusCode, err)

	return false, err
}

// HandleError is the Error handler
// for the whole tfacon.
func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
