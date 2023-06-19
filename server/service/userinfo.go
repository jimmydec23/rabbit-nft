package service

import (
	"rabbitnft/server/message"
	"rabbitnft/server/model"
)

type UserInfo struct {
	Private model.Account       `json:"-"`
	Basic   message.AccountInfo `json:"basic"`
}
