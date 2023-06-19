package message

import "math/big"

type Account struct {
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type AccountInfo struct {
	Account  string   `json:"account"`
	Nickname string   `json:"nickname"`
	Address  string   `json:"address"`
	Balance  *big.Int `json:"balance"`
}
