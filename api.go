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

var (
	paramsError = errors.New("Params Error")
	authError   = errors.New("OAuth Fail")
)

const (
	APIServer = ""
)

type API struct {
	HttpClient *http.Client
	Timeout    time.Duration
	APIServer  string
	OA         *OAuth
}

func NewAPI(ClientID string, customHttpClient *http.Client, opts ...string) (*API, error) {
	if ClientID == "" {
		return nil, paramsError
	}
	_api := &API{}
	if customHttpClient == nil {
		_api.HttpClient = http.DefaultClient
	} else {
		_api.HttpClient = customHttpClient
	}

	_api.Timeout = 30 * time.Second
	_api.APIServer = APIServer
	_api.OA = &OAuth{}
	return _api, nil

}

//Usage
// https://www.interactiveplus.org <- API Server
// /authcode <- QueryString
// Output: https://www.interactiveplus.org/authcode
func (a *API) GetFormatURL(QueryString string) string {
	return fmt.Sprintf("%s%s", a.APIServer, QueryString)
}

//access the url using GET Method
func (a *API) GetURL(URL string) (string, string, error) {
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
	return string(body), res.Status, nil
}

//access the url using POST method
// Return Value : Response Body, HTTP Code, Error
func (a *API) PostURL(URL string, Value url.Values) (string, string, error) {
	req, err := http.NewRequest("POST", a.GetFormatURL(URL), bytes.NewBufferString(Value.Encode()))
	if err != nil {
		return "", "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), a.Timeout)
	defer cancel()
	res, err := a.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return string(body), res.Status, nil
}
