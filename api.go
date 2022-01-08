package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var ()

type API struct {
	HttpClient *http.Client
	Timeout    time.Duration
	APIServer  string
}

func NewAPI(customHttpClient *http.Client, opts ...string) (*API, error) {
	if len(opts) == 0 {
		return nil, errors.New("")
	}
	_api := &API{}
	if customHttpClient == nil {
		_api.HttpClient = http.DefaultClient
	} else {
		_api.HttpClient = customHttpClient
	}

	_api.Timeout = 30 * time.Second
	_api.APIServer = ""

	return _api, nil

}

func (a *API) GetFormatURL(QueryString string) string {
	return fmt.Sprintf("%s/%s", a.APIServer, QueryString)
}

//access the url using GET Method
func (a *API) GetURL(URL string) (string, error) {
	req, err := http.NewRequest("GET", a.GetFormatURL(URL), nil)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()
	res, err := a.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body), nil
}

//access the url using POST method
func (a *API) PostURL(URL string, Value url.Values) (string, error) {
	req, err := http.NewRequest("POST", a.GetFormatURL(URL), bytes.NewBufferString(Value.Encode()))
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()
	res, err := a.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body), nil
}
