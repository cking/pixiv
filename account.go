package pixiv

import "encoding/json"

// AccountProfileImages for account
type AccountProfileImages struct {
	Px16  string `json:"px_16x16"`
	Px50  string `json:"px_50x50"`
	Px170 string `json:"px_170x170"`
}

// Account info
type Account struct {
	ID               json.Number `json:"id"`
	Name             string      `json:"name"`
	Account          string      `json:"account"`
	MailAddress      string      `json:"mail_address"`
	IsPremium        bool        `json:"is_premium"`
	XRestrict        int         `json:"x_restrict"`
	IsMailAuthorized bool        `json:"is_mail_authorized"`

	ProfileImage AccountProfileImages `json:"profile_image_urls"`
}
