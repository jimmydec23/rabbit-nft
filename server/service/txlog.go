package service

import (
	"rabbitnft/server/model"
	"rabbitnft/server/sdb"
)

type TxLogSvc struct {
	sdb *sdb.Client
}

func NewTxLogSvc(sdb *sdb.Client) *TxLogSvc {
	return &TxLogSvc{sdb: sdb}
}

func (t *TxLogSvc) TxLogList(address string, page, limit int) (int64, []model.TxLog, error) {
	return t.sdb.TxLogList(address, page, limit)
}
