package oauth

import (
	"net/url"

	"github.com/InteractivePlus/InteractiveSSO-Go/api"
	"github.com/InteractivePlus/InteractiveSSO-Go/common"
	"github.com/InteractivePlus/InteractiveSSO-Go/user"
)

type OAuthToken struct {
	AccessToken    string   `json:"access_token"`
	RefreshToken   string   `json:"refresh_token, omitempty"`
	ObtainedMethod int      `json:"obtained_method"`
	Issued         int      `json:"issued"`
	Expires        int      `json:"expires"`
	LastRenewed    int      `json:"last_renewed"`
	RefreshExpires int      `json:"refresh_expires"`
	MaskID         string   `json:"mask_id"`
	ClientID       string   `json:"client_id"`
	Scope          []string `json:"scope"`
}

type OAuthScope struct {
	Info          string `json:"info"`
	Notifications string `json:"notifications"`
	ContactSales  string `json:"contact_sales"`
}

type OAuthUserInfo struct {
	MaskID      string                  `json:"mask_id"`
	DisplayName string                  `json:"display_name"`
	Settings    *user.UserSettingEntity `json:"settings"`
}
type OAuth struct {
	api      *api.API
	Token    *OAuthToken
	Scope    *OAuthScope
	UserInfo *OAuthUserInfo
	AuthCode string
}

var (
	HTTP200OK      = "200 OK"
	HTTP201CREATED = "201 CREATED"
)

//Optional Params: code_challenge code_challenge_type	state
func (o *OAuth) GetAuthCode(UID, access_token, mask_id, client_id, scope, redirect_uri string, opts ...string) (string, error) {

}

//Optional Params: client_secret code_verifier
func (o *OAuth) GetAccessToken(isPKCE bool, clientSecret string, opts ...string) (*OAuthToken, *common.JSONError, error) {
	if o.AuthCode == "" || o.Token.ClientID == "" {
		return nil, nil, common.ParamsError
	}

	payload := &url.Values{}
	payload.Add("code", o.AuthCode)
	payload.Add("client_id", o.Token.ClientID)

	if !isPKCE {
		payload.Add("client_secret", clientSecret)
	} else {
		//PKCE Mode	Ignore ClientSecret
		codeVerifier := opts[0]
		payload.Add("code_verifier", codeVerifier)
	}

	res, status, err := o.api.PostURL("/oauth_token", payload)
	if err != nil {
		return nil, nil, err
	}

	if status != HTTP201CREATED {
		return nil, nil, common.AuthError
	}
	var ret OAuthToken
	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	o.Token = &ret
	return &ret, nil, nil

}

//Optional Params: client_secret mask_id
func (o *OAuth) VerifyAccessToken(opts ...string) (*OAuthToken, *common.JSONError, error) {
	if o.Token == nil || o.Token.ClientID == "" {
		return nil, nil, common.ParamsError
	}
	var params = map[string]string{}
	params["access_token"] = o.Token.AccessToken
	params["client_id"] = o.Token.ClientID
	if len(opts) == 1 {
		params["client_secret"] = opts[0]
	} else if len(opts) > 1 {
		params["client_secret"] = opts[0]
		params["mask_id"] = opts[1]
	}
	res, status, err := o.api.GetURLWithParams("/oauth_token/verified_status", params)
	if err != nil {
		return nil, nil, err
	}
	if status != HTTP200OK {
		return nil, nil, common.AuthError
	}

	var ret OAuthToken

	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

//Optional Params: client_secret
func (o *OAuth) RefreshAccessToken(opts ...string) (*OAuthToken, *common.JSONError, error) {
	if o.Token == nil || o.Token.ClientID == "" {
		return "", common.ParamsError
	}
	var params = map[string]string{}
	params["client_id"] = o.Token.ClientID
	if o.Token.RefreshToken != "" {
		params["refresh_token"] = o.Token.RefreshToken
	}

	if len(opts) > 0 {
		params["client_secret"] = opts[0]
	}

	res, status, err := o.api.GetURLWithParams("/oauth_token/refresh_result", params)
	if err != nil {
		return "", err
	}

	if status != HTTP200OK {
		return "", err
	}

	var ret OAuthToken

	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return res, nil

}

func (o *OAuth) GetUserInfo()
