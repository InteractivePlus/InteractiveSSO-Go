package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/InteractivePlus/InteractiveSSO-Go/oauth"
	"github.com/InteractivePlus/InteractiveSSO-Go/user"
)

const (
	APIServer = ""
)

type API struct {
	Ctx        context.Context
	HttpClient *http.Client
	Timeout    time.Duration
	APIServer  string
	o          *oauth.OAuth
	u          *user.User
}

//Usage
// https://www.interactiveplus.org <- API Server
// /authcode <- QueryString
// Output: https://www.interactiveplus.org/authcode
func (a *API) GetFormatURL(QueryString string) string {
	return fmt.Sprintf("%s%s", a.APIServer, QueryString)
}

func (a *API) ParseURLWithParams(URL string, params map[string]string) string {
	genStringSlice := []string{}
	var genString string

	for k, v := range params {
		genString = k + "=" + v
		genStringSlice = append(genStringSlice, genString)
	}
	return fmt.Sprintf("%s%s?%s", a.APIServer, URL, strings.Join(genStringSlice, "&"))
}

//access the url using GET Method
func (a *API) GetURL(URL string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", a.GetFormatURL(URL), nil)
	if err != nil {
		return nil, "", err
	}
	ctx, cancel := context.WithTimeout(a.Ctx, a.Timeout)
	defer cancel()
	res, err := a.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body, res.Status, nil
}

func (a *API) GetURLWithParams(URL string, params map[string]string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", a.ParseURLWithParams(URL, params), nil)
	if err != nil {
		return nil, "", err
	}
	ctx, cancel := context.WithTimeout(a.Ctx, a.Timeout)
	defer cancel()
	res, err := a.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body, res.Status, nil
}

//access the url using POST method
// Return Value : Response Body, HTTP Code, Error
func (a *API) PostURL(URL string, Value map[string]string) ([]byte, string, error) {
	payload, _ := json.Marshal(Value)
	req, err := http.NewRequest("POST", a.GetFormatURL(URL), bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, "", err
	}
	ctx, cancel := context.WithTimeout(a.Ctx, a.Timeout)
	defer cancel()
	res, err := a.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body, res.Status, nil
}

func (a *API) PatchURL(URL string, Value map[string]string) ([]byte, string, error) {

}

func (a *API) OAuth(ClientID string) *oauth.OAuth {
	if a.o == nil {
		//Cache it for the first time
		a.o = &oauth.OAuth{
			Token: &oauth.OAuthToken{
				ClientID: ClientID,
			},
			API: a,
		}
	}
	return a.o
}

func (a *API) User() *user.User {
	if a.u == nil {
		//Cache it for the first time
		a.u = &user.User{
			API: a,
		}
	}
	return a.u
}
