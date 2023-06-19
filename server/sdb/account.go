package sdb

import (
	"fmt"
	"rabbitnft/server/model"
)

func (c *Client) AccountRegister(account *model.Account) error {
	has, err := c.engine.Where("account = ?", account.Account).Get(&model.Account{})
	if err != nil {
		return err
	}
	if has {
		return fmt.Errorf("Account %s already exist.", account.Account)
	}
	_, err = c.engine.InsertOne(account)
	return err
}

func (c *Client) AccountDelete(account *model.Account) error {
	_, err := c.engine.Delete(account)
	return err
}

func (c *Client) AccountValidate(account *model.Account) error {
	has, err := c.engine.Where("account = ?", account.Account).
		And("password = ?", account.Password).
		Get(&model.Account{})
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf("Account %s validate failed.", account.Account)
	}
	return nil
}

func (c *Client) AccountGetByAddress(address *string) (*model.Account, error) {
	accs := []model.Account{}
	err := c.engine.Where("address = ?", address).Find(&accs)
	if err != nil {
		return nil, err
	}
	if len(accs) == 0 {
		return nil, fmt.Errorf("Account not found")
	}
	return &accs[0], nil
}

func (c *Client) AccountGet(acc *model.Account) (*model.Account, error) {
	dba := &model.Account{}
	has, err := c.engine.Where("account = ?", acc.Account).
		And("password = ? ", acc.Password).
		Get(dba)
	if !has {
		return nil, fmt.Errorf("Username or password error.")
	}
	return dba, err
}
