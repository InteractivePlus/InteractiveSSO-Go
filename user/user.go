package user

import (
	"fmt"
	"strconv"

	"github.com/InteractivePlus/InteractiveSSO-Go/common"
)

//Be careful! Bool in Go is a type not int
type UserSettingEntity struct {
	AllowEmailNotifications bool `json:"allowEmailNotifications"`
	AllowSaleEmail          bool `json:"allowSaleEmail"`
	AllowSMSNotifications   bool `json:"allowSMSNotifications"`
	AllowSaleSMS            bool `json:"allowSaleSMS"`
	AllowCallNotifications  bool `json:"allowCallNotifications"`
	AllowSaleCall           bool `json:"allowSaleCall"`
}

type MaskID struct {
	MaskId      string            `json:"mask_id"`
	ClientID    string            `json:"client_id"`
	UID         int               `json:"uid"`
	DisplayName string            `json:"display_name"`
	CreateTime  int               `json:"create_time"`
	Settings    UserSettingEntity `json:"settings"`
}

type UserEntity struct {
	UID           string            `json:"uid"`
	Username      string            `json:"username"`
	Nickname      string            `json:"nickname, omitempty"`
	Signature     string            `json:"signature, omitempty"`
	Email         string            `json:"email, omitempty"`
	Phone         string            `json:"phone, omitempty"`
	EmailVerified bool              `json:"emailVerified"`
	PhoneVerified bool              `json:"phoneVerified"`
	AccountFrozen bool              `json:"accountFrozen"`
	Settings      UserSettingEntity `json:"settings"`
}

type User struct {
	API           *api.API
	SettingEntity *UserSettingEntity
	Maskid        *MaskID
}

type RegisterRes struct {
	UID                         string `json:"uid"`
	Username                    string `json:"username"`
	Email                       string `json:"email, omitempty"`
	Phone                       string `json:"phone, omitempty"`
	PhoneVerificationSentMethod int    `json:"phoneVerificationSentMethod"`
}

type EmailRes struct {
	Username string `json:"username"`
	Nickname string `json:"nickname, omitempty"`
	Email    string `json:"email"`
}

type PhoneRes struct {
}

//Opts: email phone
func (u *User) Register(Username, Password, Captcha_id string, opts ...string) (*RegisterRes, *common.JSONError, error) {
	var params = map[string]string{}
	params["username"] = Username
	params["password"] = Password
	params["captcha_id"] = Captcha_id
	if len(opts) == 1 {
		params["email"] = opts[0]
	} else if len(opts) == 2 {
		params["email"] = opts[0]
		params["phone"] = opts[1]
	}
	res, status, err := u.API.PostURL("/user", params)
	if err != nil {
		return nil, nil, err
	}

	if status != common.HTTP201CREATED {
		return nil, nil, err
	}

	var ret RegisterRes

	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil

}

func (u *User) VerifyEmail(VeriCode string) (*EmailRes, *common.JSONError, error) {
	res, status, err := u.API.GetURL(fmt.Sprintf("%s/%s", "/vericodes/verifyPhoneResult", VeriCode))
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP200OK {
		return nil, nil, err
	}

	var ret EmailRes

	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

func (u *User) RequestEmailResend(Email, Captcha_id string) (*common.JSONError, error) {
	var params = map[string]string{}
	params["email"] = Email
	params["captcha_id"] = Captcha_id
	res, status, err := u.API.PostURL("/vericodes/sendAnotherVerifyEmailRequest", params)
	if err != nil {
		return nil, err
	}
	if status != common.HTTP201CREATED {
		return nil, err
	}
	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil
}

func (u *User) RequestPhoneResend(Preferred_send_method int, Phone, Captcha_id string) (*common.JSONError, error) {
	var params = map[string]string{}
	params["phone"] = Phone
	params["preferred_send_method"] = strconv.Itoa(Preferred_send_method)
	params["captcha_id"] = Captcha_id
	res, status, err := u.API.PostURL("/vericodes/sendAnotherVerifyEmailRequest", params)
	if err != nil {
		return nil, err
	}
	if status != common.HTTP201CREATED {
		return nil, err
	}
	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil
}
