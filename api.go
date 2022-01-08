package interactivesso

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/InteractivePlus/InteractiveSSO-Go/oauth"
	"github.com/InteractivePlus/InteractiveSSO-Go/user"
)

var (
	ParamsError      = errors.New("Params Error")
	AuthError        = errors.New("OAuth Fail")
	IsDebug     bool = false
)

const (
	NOT_SENT = itoa
	EMAIL
	SMS_MESSAGE
	PHONE_CALL
)

//Error Enum
const (
	NO_ERROR = itoa
	UNKNOWN_INNER_ERROR
	STORAGE_ENGINE_ERROR
	INNER_ARGUMENT_ERROR
	SENDER_SERVICE_ERROR
)

const (
	ITEM_NOT_FOUND_ERROR = itoa + 10
	ITEM_ALREADY_EXIST_ERROR
	ITEM_EXPIRED_OR_USED_ERROR
	PERMISSION_DENIED
	CREDENTIAL_NOT_MATCH
)

const (
	REQUEST_PARAM_FORMAT_ERROR = itoa + 20
)

const (
	APIServer = ""
)

type GeneralResult struct {
	ErrCode             int             `json:"errorCode"`
	ErrorDescription    string          `json:"errorDescription,omitempty"`
	ErrorFile           string          `json:"errorFile,omitempty"`
	ErrorLine           int             `json:"errorLine,omitempty"`
	ErrorParam          string          `json:"errorParam,omitempty"`
	Item                string          `json:"item,omitempty"`
	Credential          string          `json:"credential,omitempty"`
	UserDefinedRootData string          `json:"user-defined-root-data,omitempty"`
	Data                json.RawMessage `json:"data,omitempty"`
}

type JSONError struct {
	SpecialError     string
	ErrorDescription string
	ErrorFile        string
	ErrorLine        int
}

type API struct {
	ctx        context.Context
	HttpClient *http.Client
	Timeout    time.Duration
	APIServer  string
}

func NewAPI(ctx context.Context, customHttpClient *http.Client) *API {
	_api := &API{}
	if customHttpClient == nil {
		_api.HttpClient = http.DefaultClient
	} else {
		_api.HttpClient = customHttpClient
	}

	if ctx == nil {
		_api.ctx = context.Background()
	} else {
		_api.ctx = ctx
	}

	_api.Timeout = 30 * time.Second
	_api.APIServer = APIServer
	return _api
}

func ProcessResult(JSON []byte, cStruct interface{}) *JSONError {
	var ret GeneralResult
	if err := json.Unmarshal(JSON, &ret); err != nil {
		return &JSONError{
			ErrorDescription: err.Error(),
		}
	}

	if ret.ErrCode != NO_ERROR {
		switch ret.ErrCode {
		case INNER_ARGUMENT_ERROR, REQUEST_PARAM_FORMAT_ERROR:
			if IsDebug {
				return &JSONError{
					ErrorDescription: ret.ErrorDescription,
					SpecialError:     ret.ErrorParam,
					ErrorFile:        ret.ErrorFile,
					ErrorLine:        ret.ErrorLine,
				}
			}
			return &JSONError{
				ErrorDescription: ret.ErrorDescription,
				SpecialError:     ret.ErrorParam,
			}
		case ITEM_NOT_FOUND_ERROR, ITEM_ALREADY_EXIST_ERROR, ITEM_EXPIRED_OR_USED_ERROR:
			if IsDebug {
				return &JSONError{
					ErrorDescription: ret.ErrorDescription,
					SpecialError:     ret.Item,
					ErrorFile:        ret.ErrorFile,
					ErrorLine:        ret.ErrorLine,
				}
			}
			return &JSONError{
				ErrorDescription: ret.ErrorDescription,
				SpecialError:     ret.Item,
			}
		case CREDENTIAL_NOT_MATCH:
			if IsDebug {
				return &JSONError{
					ErrorDescription: ret.ErrorDescription,
					SpecialError:     ret.Credential,
					ErrorFile:        ret.ErrorFile,
					ErrorLine:        ret.ErrorLine,
				}
			}
			return &JSONError{
				ErrorDescription: ret.ErrorDescription,
				SpecialError:     ret.Credential,
			}
		case PERMISSION_DENIED, SENDER_SERVICE_ERROR, STORAGE_ENGINE_ERROR, UNKNOWN_INNER_ERROR:
			if IsDebug {
				return &JSONError{
					ErrorDescription: ret.ErrorDescription,
					ErrorFile:        ret.ErrorFile,
					ErrorLine:        ret.ErrorLine,
				}
			}
			return &JSONError{
				ErrorDescription: ret.ErrorDescription,
			}
		default:
			return &JSONError{
				ErrorDescription: "Unknown Error",
			}
		}
	}
	if err := json.Unmarshal(ret.Data, cStruct); err != nil {
		return &JSONError{
			ErrorDescription: err.Error(),
		}
	}
	return nil
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
	ctx, cancel := context.WithTimeout(a.ctx, a.Timeout)
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
	ctx, cancel := context.WithTimeout(a.ctx, a.Timeout)
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
func (a *API) PostURL(URL string, Value url.Values) ([]byte, string, error) {
	req, err := http.NewRequest("POST", a.GetFormatURL(URL), bytes.NewBufferString(Value.Encode()))
	if err != nil {
		return nil, "", err
	}
	ctx, cancel := context.WithTimeout(a.ctx, a.Timeout)
	defer cancel()
	res, err := a.HttpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body, res.Status, nil
}

func (a *API) OAuth(ClientID string) *oauth.OAuth {
	return &oauth.OAuth{
		Token: &oauth.OAuthToken{
			ClientID: ClientID,
		},
		api: a,
	}
}

func (a *API) User() *user.User {
	return &user.User{
		api: a,
	}
}
