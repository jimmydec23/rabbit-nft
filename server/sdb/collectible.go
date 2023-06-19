package sdb

import (
	"fmt"
	"rabbitnft/server/model"
)

func (c *Client) CollectibleAdd(collectible *model.Collectible) error {
	_, err := c.engine.InsertOne(collectible)
	return err
}

func (c *Client) CollectibleUpdate(collectible *model.Collectible) error {
	i, err := c.engine.Where("token_id = ?", collectible.TokenId).AllCols().Update(collectible)
	if i == 0 {
		return Err_NoEffect
	}
	return err
}

func (c *Client) CollectibleGet(tokenId int64) (*model.Collectible, error) {
	collectible := &model.Collectible{TokenId: tokenId}
	_, err := c.engine.Get(collectible)
	return collectible, err
}

func (c *Client) CollectibleList(condi *model.Collectible, page, limit int) (int64, []model.Collectible, error) {
	offset := (page - 1) * limit
	list := []model.Collectible{}
	total, err := c.engine.UseBool("onsale").Limit(limit, offset).FindAndCount(&list, condi)
	return total, list, err
}

func (c *Client) CollectibleCounterGet() (int64, error) {
	counter := &model.CollectibleCounter{}
	has, err := c.engine.Get(counter)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, fmt.Errorf("Counter not exist")
	}
	return counter.Counter, nil
}

func (c *Client) CollectibleCounterUpdate(newCount int64) error {
	_, err := c.engine.Update(&model.CollectibleCounter{Counter: newCount})
	return err
}

func (c *Client) CollectibleOwnerUpdate(owner *model.CollectibleOwner) error {
	condi := &model.CollectibleOwner{TokenId: owner.TokenId}
	has, err := c.engine.Exist(condi)
	if err != nil {
		return err
	}
	if has {
		_, err := c.engine.Update(owner, condi)
		return err
	} else {
		_, err := c.engine.InsertOne(owner)
		return err
	}
}

func (c *Client) CollectibleOfOwner(owner string) ([]model.CollectibleOwner, error) {
	list := []model.CollectibleOwner{}
	err := c.engine.Find(&list, &model.CollectibleOwner{Owner: owner})
	return list, err
}
