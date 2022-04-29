package rest

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"

	"github.com/Clevermind-Think-Mint/buymint-cli-go/internal/logger"
	//cookiejar "github.com/juju/persistent-cookiejar"
)

var client http.Client
var jar *cookiejar.Jar

type customPSL struct{}

func (customPSL) String() string {
	return "customPSL"
}
func (customPSL) PublicSuffix(d string) string {
	return d[strings.LastIndex(d, ".")+1:]
}

// Get is...
func Get(URL string, headers map[string]string, options map[string]interface{}) (int, []byte, map[string]string, error) {
	// Init headers if needed
	if headers == nil {
		headers = make(map[string]string, 0)
	}
	// Fetching
	return fetch(http.MethodGet, URL, headers, nil, options)
}

// Post is...
func Post(URL string, data interface{}, headers map[string]string, options map[string]interface{}) (int, []byte, map[string]string, error) {
	// Init options if needed
	if options == nil {
		options = make(map[string]interface{}, 0)
	}
	// Init headers if needed
	if headers == nil {
		headers = make(map[string]string, 0)
	}
	var bodyRequest io.Reader
	switch data.(type) {
	case url.Values:
		if headers["Content-Type"] == "" {
			headers["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
		}
		bodyRequest = strings.NewReader(data.(url.Values).Encode())
	default:
		switch reflect.ValueOf(data).Kind() {
		case reflect.Map:
			bodyJSON, err := json.Marshal(data)
			if err != nil {
				return 0, nil, nil, err
			}
			if headers["Content-Type"] == "" {
				headers["Content-Type"] = "application/json"
			}
			bodyRequest = bytes.NewReader(bodyJSON)
		default:
			return 0, nil, nil, errors.New("Unsupported body request type/kind")
		}
	}
	// Fetching
	return fetch(http.MethodPost, URL, headers, bodyRequest, options)
}

func fetch(method string, URL string, headers map[string]string, bodyRequest io.Reader, options map[string]interface{}) (int, []byte, map[string]string, error) {
	if options == nil {
		options = make(map[string]interface{}, 0)
	}
	if options["redirect"] == nil {
		options["redirect"] = true
	}
	if options["IgnoreInsecureSsl"] == nil {
		options["IgnoreInsecureSsl"] = false
	}

	// Building request
	request, err := http.NewRequest(method, URL, bodyRequest)
	if err != nil {
		return 0, nil, nil, err
	}
	// Setting headers
	for field, value := range headers {
		request.Header.Set(field, value)
	}
	logger.Debug("Sending HTTP request:\n%s", dumpRequest(request))
	// Init cookie if needed
	if jar == nil {
		tmpJar, err := cookiejar.New(&cookiejar.Options{
			PublicSuffixList: customPSL{},
			//Filename:         "./dist/cookie.json",
		})
		jar = tmpJar
		if err != nil {
			return 0, nil, nil, err
		}
	}
	originDomain := ""
	response, err := (&http.Client{
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: options["IgnoreInsecureSsl"].(bool)},
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				conn, err := net.Dial(network, addr)
				if err == nil {
					originDomain = conn.LocalAddr().String()
				}
				return conn, err
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !options["redirect"].(bool) {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}).Do(request)
	if err != nil {
		return 0, nil, nil, err
	}
	logger.Debug("Received HTTP response:\n%s", dumpResponse(response))
	// Reading headers
	headers = make(map[string]string, 0)
	for field, values := range response.Header {
		for _, value := range values {
			headers[field] = value
		}
	}
	// TODO: add options to choose this behaviour
	// Forcing Origin
	if headers["Origin"] == "" {
		headers["Origin"] = originDomain
	}
	// Reading body
	defer response.Body.Close()
	bodyResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, nil, headers, err
	}
	// Checking status code
	if response.StatusCode >= 400 {
		return response.StatusCode, bodyResponse, headers, fmt.Errorf("Status code: %d", response.StatusCode)
	}
	// Returning body with no error
	return response.StatusCode, bodyResponse, headers, nil
}

func dumpRequest(request *http.Request) string {
	dump, _ := httputil.DumpRequest(request, true)
	return string(dump)
}

// dumpResponse is...
func dumpResponse(request *http.Response) string {
	dump, _ := httputil.DumpResponse(request, true)
	return string(dump)
}
