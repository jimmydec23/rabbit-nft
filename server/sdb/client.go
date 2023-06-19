package sdb

import (
	"database/sql"
	"rabbitnft/server/log"
	"rabbitnft/server/model"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

type Client struct {
	db     *sql.DB
	engine *xorm.Engine
}

func NewClient() (*Client, error) {
	dbpath := "db.sqlite3"
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	engine, err := xorm.NewEngine("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}

	c := &Client{
		db:     db,
		engine: engine,
	}
	return c, nil
}

func (c *Client) InitTable() error {
	// init account table
	accountTbl := &model.Account{}
	exist, err := c.engine.IsTableExist(accountTbl)
	if err != nil {
		return err
	}
	if !exist {
		err := c.engine.CreateTables(accountTbl)
		if err != nil {
			return err
		}
		log.Logger.Info("Account table created.")
	}

	// init collectible table
	colTable := &model.Collectible{}
	exist, err = c.engine.IsTableExist(colTable)
	if err != nil {
		return err
	}
	if !exist {
		err := c.engine.CreateTables(colTable)
		if err != nil {
			return err
		}
		log.Logger.Info("Collectible table created.")
	}

	// init collectible counter
	counterTbl := &model.CollectibleCounter{}
	exist, err = c.engine.IsTableExist(counterTbl)
	if err != nil {
		return err
	}
	if !exist {
		err := c.engine.CreateTables(counterTbl)
		if err != nil {
			return err
		}
		_, err = c.engine.InsertOne(&model.CollectibleCounter{Counter: CounterStart})
		if err != nil {
			return err
		}
		log.Logger.Info("Collectible counter table created.")
	}

	// init collectible owner
	ownerTbl := &model.CollectibleOwner{}
	exist, err = c.engine.IsTableExist(ownerTbl)
	if err != nil {
		return err
	}
	if !exist {
		err := c.engine.CreateTables(ownerTbl)
		if err != nil {
			return err
		}
		log.Logger.Info("Collectible owner table created.")
	}

	// init tx log table
	txlogTbl := &model.TxLog{}
	exist, err = c.engine.IsTableExist(txlogTbl)
	if err != nil {
		return err
	}
	if !exist {
		err := c.engine.CreateTables(txlogTbl)
		if err != nil {
			return err
		}
		log.Logger.Info("TxLog table created.")
	}

	log.Logger.Info("Table init success.")
	return nil
}

func (c *Client) Close() {
	c.db.Close()
}
