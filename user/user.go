package user


import (
	"github.com/InteractivePlus/InteractiveSSO-Go/common"
)

//Be careful! Bool in Go is a type not int
type UserSettingEntity struct {
	AllowEmailNotifications bool 	`json:"allowEmailNotifications"`
	AllowSaleEmail			bool	`json:"allowSaleEmail"`
	AllowSMSNotifications	bool	`json:"allowSMSNotifications"`
	AllowSaleSMS			bool	`json:"allowSaleSMS"`
	AllowCallNotifications	bool	`json:"allowCallNotifications"`
	AllowSaleCall			bool	`json:"allowSaleCall"`
}

type MaskID struct {
	MaskId 		 string		`json:"mask_id"`
	ClientID	 string		`json:"client_id"`
	UID			 int     	`json:"uid"`
	DisplayName	 string		`json:"display_name"`
	CreateTime	 int		`json:"create_time"`
	Settings	*UserSettingEntity	`json:"settings"`
}
type User struct {
	api           *interactivesso.API
	SettingEntity *UserSettingEntity
	Maskid		  *MaskID
}

func (u *User) 

