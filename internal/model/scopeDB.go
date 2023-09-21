package model

import (
	conf "bittoCralwer/config"
	"bittoCralwer/internal/dao"
	"bittoCralwer/internal/protocol/dto"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// ScopeDB : 블록들 정보
type ScopeDB struct {
	DB    *sql.DB // Add this line
	start chan struct{}
}

// NewScopeDB : ScopeDB 객체 할당 및 반환
func NewScopeDB(config *conf.Config, root *Repositories) (IRepository, error) {
	cfg := config.Repositories["scope-db"]
	r := &ScopeDB{
		start: make(chan struct{}),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg["username"], cfg["pass"], cfg["datasource"], cfg["db"])
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		// TODO: log 통일
		log.Fatal(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		db.Close()
		return nil, err
	}

	r.DB = db
	// TODO: trace db connection log

	return r, nil
}

func (p *ScopeDB) Close() {
	p.DB.Close()
}

func (p *ScopeDB) Start() error {
	return func() (err error) {
		defer func() {
			if v := recover(); v != nil {
				err = v.(error)
			}
		}()
		close(p.start)
		return
	}()
}

func (p *ScopeDB) SaveEthereumBlock(block *dto.EthereumBlock) error {
	blockDAO := dao.NewEthereumBlockDAO(p.DB)
	return blockDAO.Save(block)
}

func (p *ScopeDB) GetEthereumBlock(blockHash string) (*dto.EthereumBlock, error) {
	blockDAO := dao.NewEthereumBlockDAO(p.DB)
	return blockDAO.Get(blockHash)
}
