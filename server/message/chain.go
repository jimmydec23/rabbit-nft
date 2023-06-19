package message

type ChainInfo struct {
	NetworkId uint64 `json:"networkId"`
	BlockNum  uint64 `json:"blockNum"`
	PeerCount uint64 `json:"peerCount"`
}

type BlockNumber struct {
	BlockNum uint64 `json:"blockNum"`
}

type BlockNumberRange struct {
	From uint64 `json:"from"`
	To   uint64 `json:"to"`
}

type Block struct {
	Hash        string        `json:"hash"`
	Txs         []Transaction `json:"txs"`
	Uncles      int           `json:"uncles"`
	ParentHash  string        `json:"parentHash"`
	UncleHash   string        `json:"sha3Uncles"`
	Coinbase    string        `json:"miner"`
	Root        string        `json:"stateRoot"`
	TxHash      string        `json:"transactionsRoot"`
	ReceiptHash string        `json:"receiptsRoot"`
	Difficulty  string        `json:"difficulty"`
	Number      uint64        `json:"number"`
	GasLimit    uint64        `json:"gasLimit"`
	GasUsed     uint64        `json:"gasUsed"`
	Time        uint64        `json:"timestamp"`
	MixDigest   string        `json:"mixHash"`
	Nonce       uint64        `json:"nonce"`
}

type Transaction struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Cost     string `json:"cost"`
	GasUsage uint64 `json:"gasUsage"`
}

type TransactionHash struct {
	Hash string `json:"hash"`
}

type BlockHash struct {
	Hash string `json:"hash"`
}
