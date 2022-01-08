package interactivesso

import (
	"context"
	"net/http"
	"time"

	"github.com/InteractivePlus/InteractiveSSO-Go/api"
)

func NewAPI(ctx context.Context, customHttpClient *http.Client) *api.API {
	_api := &api.API{}
	if customHttpClient == nil {
		_api.HttpClient = http.DefaultClient
	} else {
		_api.HttpClient = customHttpClient
	}

	if ctx == nil {
		_api.Ctx = context.Background()
	} else {
		_api.Ctx = ctx
	}

	_api.Timeout = 30 * time.Second
	_api.APIServer = api.APIServer
	return _api
}
