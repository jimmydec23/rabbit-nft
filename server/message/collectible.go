package message

type TokenId struct {
	TokenId int64 `json:"tokenId"`
}

type Collectible struct {
	TokenId     int64  `json:"tokenId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Hash        string `json:"hash"`
	Url         string `json:"url"`
	Owner       string `json:"owner"`
	Onsale      bool   `json:"onsale"`
	Price       string `json:"price"`
	Approver    string `json:"approver"`
}

type MintCollectible struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Hash        string `json:"hash"`
}

type SellCollectible struct {
	TokenId int64  `json:"tokenId"`
	Price   string `json:"price"`
}
