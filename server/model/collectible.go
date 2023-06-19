package model

import "encoding/json"

type Collectible struct {
	TokenId     int64  `json:"tokenId" xorm:"token_id"`
	Name        string `json:"name"`
	Hash        string `json:"hash"`
	Description string `json:"description"`
	Onsale      bool   `json:"onsale"`
	Price       string `json:"price"`
}

func (c *Collectible) Bytes() []byte {
	b, _ := json.Marshal(c)
	return b
}

type CollectibleCounter struct {
	Counter int64 `json:"counter"`
}

type CollectibleOwner struct {
	TokenId int64  `json:"tokenId" xorm:"token_id"`
	Owner   string `json:"owner"`
}
