package common

import (
	"encoding/json"
	"errors"
)

type SettingBoolean struct {
}

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

type JSONError struct {
	SpecialError     string
	ErrorDescription string
	ErrorFile        string
	ErrorLine        int
}

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
