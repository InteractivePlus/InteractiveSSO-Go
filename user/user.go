package user

import (
	"fmt"
	"strconv"

	"github.com/InteractivePlus/InteractiveSSO-Go/common"
)

const (
	EMAIL_NOT_VERIFIED = iota + 1
	PHONE_NOT_VERIFIED
	EITHER_NOT_VERIFIED
	ACCOUNT_FROZEN
	UNKNOWN
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

type MaskIDEntity struct {
	MaskId      string            `json:"mask_id"`
	ClientID    string            `json:"client_id"`
	UID         int               `json:"uid"`
	DisplayName string            `json:"display_name"`
	CreateTime  int               `json:"create_time"`
	Settings    UserSettingEntity `json:"settings"`
}

type UserEntity struct {
	UID           int               `json:"uid"`
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
	API *api.API
}

type RegisterRes struct {
	UID                         int    `json:"uid"`
	Username                    string `json:"username"`
	Email                       string `json:"email, omitempty"`
	Phone                       string `json:"phone, omitempty"`
	PhoneVerificationSentMethod int    `json:"phoneVerificationSentMethod"`
}

type VerifyEmailRes struct {
	Username string `json:"username"`
	Nickname string `json:"nickname, omitempty"`
	Email    string `json:"email"`
}

type VerifyPhoneRes struct {
	Username string `json:"username"`
	Nickname string `json:"nickname, omitempty"`
	Phone    string `json:"phone"`
}

type LoginRes struct {
	AccessToken   string     `json:"access_token, omitempty"`
	RefreshToken  string     `json:"refresh_token, omitempty"`
	ExpireTime    int        `json:"expire_time, omitempty"`
	RefreshExpire int        `json:"refresh_expire, omitempty"`
	User          UserEntity `json:"user, omitempty"`
	ErrorReason   int        `json:"errorReason, omitempty"`
	Email         string     `json:"email, omitempty"`
	Phone         string     `json:"phone, omitempty"`
	UID           int        `json:"uid, omitempty"`
}

type ModifyUserPayload struct {
	UID         int               `json:"uid"`
	AccessToken string            `json:"access_token"`
	Nickname    string            `json:"nickname, omitempty"`
	Signature   string            `json:"signature, omitempty"`
	Settings    UserSettingEntity `json:"settings"`
}

type MaskPayload struct {
	UID         int               `json:"uid"`
	AccessToken string            `json:"access_token"`
	ClientID    string            `json:"client_id, omitempty"`
	DisplayName string            `json:"display_name, omitempty"`
	Settings    UserSettingEntity `json:"settings, omitempty"`
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

func (u *User) VerifyEmail(VeriCode string) (*VerifyEmailRes, *common.JSONError, error) {
	res, status, err := u.API.GetURL(fmt.Sprintf("/vericodes/verifyEmailResult/%s", VeriCode))
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP200OK {
		return nil, nil, err
	}

	var ret VerifyEmailRes

	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

func (u *User) VerifyPhone(UID int, VeriCode string) (*VerifyPhoneRes, *common.JSONError, error) {
	res, status, err := u.API.GetURL(fmt.Sprintf("/vericodes/verifyPhoneResult/%s?uid=%d", VeriCode, UID))
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP200OK {
		return nil, nil, err
	}

	var ret VerifyPhoneRes

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

func (u *User) RequestPhoneResend(Preferred_send_method int, Phone, Captcha_id string) (*common.SENT_METHOD, *common.JSONError, error) {
	var params = map[string]string{}
	params["phone"] = Phone
	params["preferred_send_method"] = strconv.Itoa(Preferred_send_method)
	params["captcha_id"] = Captcha_id
	res, status, err := u.API.PostURL("/vericodes/sendAnotherVerifyEmailRequest", params)
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP201CREATED {
		return nil, nil, err
	}

	var ret common.SENT_METHOD
	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

//Opts: Username, Phone, Email
//Leave it ""
func (u *User) Login(Password, Captcha_id, Username, Phone, Email string) (*LoginRes, *common.JSONError, error) {
	if Password == "" || Captcha_id == "" {
		return nil, nil, common.ParamsError
	}
	var params = map[string]string{}
	if Username != "" {
		params["username"] = Username
	}
	if Phone != "" {
		params["phone"] = Phone
	}

	if Email != "" {
		params["email"] = Email
	}

	res, status, err := u.API.PostURL("/user/token", params)
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP201CREATED {
		return nil, nil, err
	}

	var ret LoginRes
	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil

}

func (u *User) VerifyToken(UID int, AccessToken string) (*common.JSONError, error) {
	if AccessToken == "" {
		return nil, common.ParamsError
	}

	res, status, err := u.API.GetURL(fmt.Sprintf("/user/%d/token/%s/checkTokenResult", UID, AccessToken))
	if err != nil {
		return nil, err
	}
	if status != common.HTTP200OK {
		return nil, err
	}

	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil
}

func (u *User) RefreshLoginInfo(UID int, RefreshToken string) (*LoginRes, *common.JSONError, error) {
	if RefreshToken == "" {
		return nil, nil, common.ParamsError
	}

	res, status, err := u.API.GetURL(fmt.Sprintf("/user/%d/token/refreshResult?refresh_token=%s", UID, RefreshToken))
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP201CREATED {
		return nil, nil, err
	}

	var ret LoginRes

	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

func (u *User) Logout(UID int, AccessToken string) (*common.JSONError, error) {
	if AccessToken == "" {
		return nil, common.ParamsError
	}

	res, status, err := u.API.DeleteURL(fmt.Sprintf("/user/%d/token/%s", UID, AccessToken))
	if err != nil {
		return nil, err
	}
	if status != common.HTTP204NOCONTENT {
		return nil, err
	}

	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil

}

func (u *User) RequestEmailVeriCode(UID int, AccessToken, NewEmail string, Preferred_send_method int) (*common.SENT_METHOD, *common.JSONError, error) {
	if AccessToken == "" || NewEmail == "" {
		return nil, nil, common.ParamsError
	}

	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["preferred_send_method"] = strconv.Itoa(Preferred_send_method)
	params["new_email"] = NewEmail
	params["access_token"] = AccessToken
	res, status, err := u.API.PostURL("/vericodes/changeEmailAddrRequest", params)
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP201CREATED {
		return nil, nil, err
	}

	var ret common.SENT_METHOD
	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

func (u *User) RequestPhoneVeriCode(UID int, AccessToken, NewPhone string, Preferred_send_method int) (*common.SENT_METHOD, *common.JSONError, error) {
	if AccessToken == "" || NewPhone == "" {
		return nil, nil, common.ParamsError
	}

	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["preferred_send_method"] = strconv.Itoa(Preferred_send_method)
	params["new_email"] = NewPhone
	params["access_token"] = AccessToken
	res, status, err := u.API.PostURL("/vericodes/changePhoneNumberRequest", params)
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP201CREATED {
		return nil, nil, err
	}

	var ret common.SENT_METHOD
	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

func (u *User) AddEmail(UID int, AccessToken, NewEmail string) (*common.JSONError, error) {
	if AccessToken == "" || NewEmail == "" {
		return nil, common.ParamsError
	}

	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["new_email"] = NewEmail
	params["access_token"] = AccessToken
	res, status, err := u.API.PatchURL("/user/email", params)
	if err != nil {
		return nil, err
	}
	if status != common.HTTP200OK {
		return nil, err
	}

	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil
}

func (u *User) ModifyEmail(UID int, VeriCode string) (*common.JSONError, error) {
	if VeriCode == "" {
		return nil, common.ParamsError
	}

	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["veriCode"] = VeriCode
	res, status, err := u.API.PatchURL("/user/email", params)
	if err != nil {
		return nil, err
	}
	if status != common.HTTP200OK {
		return nil, err
	}

	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil
}

func (u *User) AddPhone(UID int, AccessToken, NewPhone string) (*common.JSONError, error) {
	if AccessToken == "" || NewPhone == "" {
		return nil, common.ParamsError
	}

	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["new_email"] = NewPhone
	params["access_token"] = AccessToken
	res, status, err := u.API.PatchURL("/user/phoneNum", params)
	if err != nil {
		return nil, err
	}
	if status != common.HTTP200OK {
		return nil, err
	}

	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil
}

func (u *User) ModifyPhone(UID int, VeriCode string) (*common.JSONError, error) {
	if VeriCode == "" {
		return nil, common.ParamsError
	}

	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["veriCode"] = VeriCode
	res, status, err := u.API.PatchURL("/user/phoneNum", params)
	if err != nil {
		return nil, err
	}
	if status != common.HTTP200OK {
		return nil, err
	}

	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil
}

func (u *User) RequestChangePasswordVeriCode(UID int, AccessToken string, Preferred_send_method int) (*common.SENT_METHOD, *common.JSONError, error) {
	if AccessToken == "" {
		return nil, nil, common.ParamsError
	}

	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["preferred_send_method"] = strconv.Itoa(Preferred_send_method)
	params["access_token"] = AccessToken
	res, status, err := u.API.PostURL("/vericodes/changePasswordRequest", params)
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP201CREATED {
		return nil, nil, err
	}

	var ret common.SENT_METHOD
	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

func (u *User) RequestResetPasswordVeriCode(UID int, AccessToken string, Preferred_send_method int, Username, Phone, Email string) (*common.SENT_METHOD, *common.JSONError, error) {
	if AccessToken == "" {
		return nil, nil, common.ParamsError
	}

	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["preferred_send_method"] = strconv.Itoa(Preferred_send_method)
	params["access_token"] = AccessToken

	if Username != "" {
		params["username"] = Username
	}
	if Phone != "" {
		params["phone"] = Phone
	}

	if Email != "" {
		params["email"] = Email
	}

	res, status, err := u.API.PostURL("/vericodes/changePasswordRequest", params)
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP201CREATED {
		return nil, nil, err
	}

	var ret common.SENT_METHOD
	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

func (u *User) ChangePassword(UID int, NewPassword, VeriCode string) (*common.JSONError, error) {
	if VeriCode == "" || NewPassword == "" {
		return nil, common.ParamsError
	}

	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["veriCode"] = VeriCode
	params["new_password"] = NewPassword
	res, status, err := u.API.PatchURL("/user/password", params)
	if err != nil {
		return nil, err
	}
	if status != common.HTTP200OK {
		return nil, err
	}

	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil
}

func (u *User) ResetPassword(UID int, NewPassword, VeriCode string, Username, Phone, Email string) (*common.JSONError, error) {
	if VeriCode == "" || NewPassword == "" {
		return nil, common.ParamsError
	}

	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["veriCode"] = VeriCode
	params["new_password"] = NewPassword

	if Username != "" {
		params["username"] = Username
	}
	if Phone != "" {
		params["phone"] = Phone
	}

	if Email != "" {
		params["email"] = Email
	}

	res, status, err := u.API.PatchURL("/user/password", params)
	if err != nil {
		return nil, err
	}
	if status != common.HTTP200OK {
		return nil, err
	}

	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil
}

func (u *User) ModifyUserInfo(UID int, AccessToken, Nickname, Signature string, Settings *UserSettingEntity) (*common.JSONError, error) {
	if AccessToken == "" {
		return nil, common.ParamsError
	}
	params := &ModifyUserPayload{
		UID:         UID,
		AccessToken: AccessToken,
	}
	if Nickname != "" {
		params.Nickname = Nickname
	}
	if Signature != "" {
		params.Signature = Signature
	}

	if Settings != nil {
		params.Settings = *Settings
	}

	res, status, err := u.API.PatchURL("/user/password", params)
	if err != nil {
		return nil, err
	}
	if status != common.HTTP200OK {
		return nil, err
	}

	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil
}

func (u *User) ListMask(UID int, AccessToken string, opts ...string) (*MaskIDEntity, *common.JSONError, error) {
	if AccessToken == "" {
		return nil, nil, common.ParamsError
	}
	var URL string
	if len(opts) > 0 {
		URL = fmt.Sprintf("/masks/%s", opts[0])
	} else {
		URL = "/masks"
	}
	var params = map[string]string{}
	params["uid"] = strconv.Itoa(UID)
	params["access_token"] = AccessToken
	res, status, err := u.API.GetURLWithParams(URL, params)
	if err != nil {
		return nil, nil, err
	}

	if status != common.HTTP200OK {
		return nil, nil, err
	}

	var ret MaskIDEntity

	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil

}

func (u *User) AddMask(UID int, AccessToken, ClientID, DisplayName string, Settings UserSettingEntity) (*MaskIDEntity, *common.JSONError, error) {
	if AccessToken == "" {
		return nil, nil, common.ParamsError
	}
	params := &MaskPayload{
		UID:         UID,
		AccessToken: AccessToken,
		DisplayName: DisplayName,
		Settings:    Settings,
	}

	res, status, err := u.API.PostURL(fmt.Sprintf("/masks/%s", ClientID), params)
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP201CREATED {
		return nil, nil, err
	}

	var ret MaskIDEntity
	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

func (u *User) ModifyMask(UID int, MaskID, AccessToken, ClientID, DisplayName string, Settings *UserSettingEntity) (*MaskIDEntity, *common.JSONError, error) {
	if AccessToken == "" {
		return nil, nil, common.ParamsError
	}
	params := &MaskPayload{
		UID:         UID,
		AccessToken: AccessToken,
		ClientID:    ClientID,
	}
	if DisplayName != "" {
		params.DisplayName = DisplayName
	}

	if Settings != nil {
		params.Settings = *Settings
	}

	res, status, err := u.API.PatchURL(fmt.Sprintf("/masks/%s", MaskID), params)
	if err != nil {
		return nil, nil, err
	}
	if status != common.HTTP200OK {
		return nil, nil, err
	}

	var ret MaskIDEntity
	if err := common.ProcessResult(res, &ret); err != nil {
		return nil, err, nil
	}

	return &ret, nil, nil
}

func (u *User) DeleteMask(UID int, MaskID, AccessToken string) (*common.JSONError, error) {
	if AccessToken == "" || MaskID == "" {
		return nil, common.ParamsError
	}

	res, status, err := u.API.DeleteURL(fmt.Sprintf("/masks/%s", MaskID))
	if err != nil {
		return nil, err
	}
	if status != common.HTTP204NOCONTENT {
		return nil, err
	}

	if err := common.FetchError(res); err != nil {
		return err, nil
	}

	return nil, nil

}
