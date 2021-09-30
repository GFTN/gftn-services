// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package gasserviceclient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/GFTN/gftn-services/utility/common"
	global_environment "github.com/GFTN/gftn-services/utility/global-environment"
)

type HTTPClientMock struct {
	GetFunc  func(url string) (*http.Response, error)
	PostFunc func(url, contentType string, body io.Reader) (*http.Response, error)
}

func (m *HTTPClientMock) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return m.PostFunc(url, contentType, body)
}
func (m *HTTPClientMock) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}

func TestSubmitTxeStatusOK(t *testing.T) {

	httpClientMock := HTTPClientMock{
		GetFunc: func(url string) (*http.Response, error) {
			return nil, nil
		},
		PostFunc: func(url, contentType string, body io.Reader) (*http.Response, error) {
			response := http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"title":"Success", "ledger":9999, "hash":"asdfn" }`)),
			}
			return &response, nil
		},
	}
	c := Client{
		HTTP: &httpClientMock,
		URL:  "",
	}
	_, _, err := c.SubmitTxe("test")
	Convey("Successful get caller identity", t, func() {
		So(err, ShouldBeNil)
	})
}

func TestSubmitTxeStatusBadRequest(t *testing.T) {

	httpClientMock := HTTPClientMock{
		GetFunc: func(url string) (*http.Response, error) {
			return nil, nil
		},
		PostFunc: func(url, contentType string, body io.Reader) (*http.Response, error) {
			response := http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"title":"Fail", "failure_reason":{"StatusBadRequest":"StatusBadRequest"} }`)),
			}
			return &response, nil
		},
	}
	c := Client{
		HTTP: &httpClientMock,
		URL:  "",
	}
	_, _, err := c.SubmitTxe("test")
	Convey("Successful get caller identity", t, func() {
		So(err, ShouldNotBeNil)
	})
}

func TestSubmitTxeStatusForbidden(t *testing.T) {

	httpClientMock := HTTPClientMock{
		GetFunc: func(url string) (*http.Response, error) {
			return nil, nil
		},
		PostFunc: func(url, contentType string, body io.Reader) (*http.Response, error) {
			response := http.Response{
				StatusCode: http.StatusForbidden,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"title":"Fail", "failure_reason":{"StatusForbidden":"StatusForbidden"} }`)),
			}
			return &response, nil
		},
	}
	c := Client{
		HTTP: &httpClientMock,
		URL:  "",
	}
	_, _, err := c.SubmitTxe("test")
	Convey("Successful get caller identity", t, func() {
		So(err, ShouldNotBeNil)
	})
}

func setEnvVariables() {
	os.Setenv(global_environment.ENV_KEY_SERVICE_NAME, "gas-service")
	os.Setenv(global_environment.ENV_KEY_HOME_DOMAIN_NAME, "ww")
	os.Setenv(global_environment.ENV_KEY_AWS_REGION, "ap-southeast-1")
	os.Setenv(global_environment.ENV_KEY_STELLAR_NETWORK, "Standalone World Wire Network ; Mar 2019")
	os.Setenv(global_environment.ENV_KEY_SERVICE_LOG_FILE, common.Abs("log.txt"))
	os.Setenv(global_environment.ENV_KEY_GAS_SVC_URL, "http://localhost:8099")
	os.Setenv(global_environment.ENV_KEY_GAS_ACCOUNT_ATTEMPTS, "3")
	os.Setenv(global_environment.ENV_KEY_WAIT_UNLOCK_DURATION, "2")
}

func TestGetAccountSequence(t *testing.T) {
	setEnvVariables()
	Convey("lock all account", t, func() {
		var fail_index int
		var success_index int
		accountNum := 30
		var mutex sync.Mutex
		wg := sync.WaitGroup{}
		wg.Add(accountNum)
		gasServiceClient := Client{
			HTTP: &http.Client{Timeout: time.Second * 20},
			URL:  os.Getenv(global_environment.ENV_KEY_GAS_SVC_URL),
		}
		for i := 0; i < accountNum; i++ {
			go func(wg *sync.WaitGroup) {
				_, _, err := gasServiceClient.GetAccountAndSequence()
				mutex.Lock()
				if err != nil {
					fail_index++
					fmt.Printf("Failed account: %v \n", fail_index)
					fmt.Println(err)
				} else {
					success_index++
					fmt.Printf("Got account: %v \n", success_index)
				}
				mutex.Unlock()
				wg.Done()
			}(&wg)
		}
		wg.Wait()
		fmt.Printf("Sent %v accounts request, got %v, and fail to get %v account \n", accountNum, success_index, fail_index)

		So(success_index+fail_index, ShouldEqual, accountNum)
	})
	/*
		Convey("sequential get lock account", t, func() {
			var fail_index int
			var success_index int
			gasServiceClient := Client{
				HTTP: &http.Client{Timeout: time.Second * 20},
				URL:  os.Getenv(global_environment.ENV_KEY_GAS_SVC_URL),
			}
			accountNum := 20
			for i := 0; i < accountNum; i++ {
				_, _, err := gasServiceClient.GetAccountAndSequence()
				if err != nil {
					fail_index++
					fmt.Printf("Failed account: %v \n", fail_index)
					fmt.Println(err)
				} else {
					success_index++
					fmt.Printf("Got account: %v \n", success_index)
				}

			}
			fmt.Printf("Sent %v accounts request, got %v, and fail to get %v account \n", accountNum, success_index, fail_index)

			So(success_index+fail_index, ShouldEqual, accountNum)
		})
	*/
}
