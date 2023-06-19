package model

const (
	TxLogType_Mint = iota
	TxLogType_Transfer
	TxLogType_Sell
	TxLogType_Purchase
)

type TxLog struct {
	Tx      string `json:"tx"`
	Type    int    `json:"type"`
	From    string `json:"from"`
	To      string `json:"to"`
	Created int64  `json:"created"`
}
