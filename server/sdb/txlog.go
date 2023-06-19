package sdb

import (
	"rabbitnft/server/model"
)

func (c *Client) TxLogRecord(r *model.TxLog) error {
	_, err := c.engine.InsertOne(r)
	return err
}

func (c *Client) TxLogList(address string, page, limit int) (int64, []model.TxLog, error) {
	list := []model.TxLog{}
	start := (page - 1) * limit
	total, err := c.engine.
		Where("`from` = ?", address).
		Or("`to` = ?", address).
		Limit(limit, start).FindAndCount(&list)
	return total, list, err
}
