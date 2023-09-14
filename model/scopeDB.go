package model

import (
	conf "bittoCralwer/config"
	"database/sql"
	"log"
)

// ScopeDB : 블록들 정보
type ScopeDB struct {
	//client         *mongo.Client
	//colContract    *mongo.Collection // Contracts
	//colExtContract *mongo.Collection // external_contracts
	//colEventManage *mongo.Collection // event_manage

	// cacheChainLock sync.RWMutex
	start chan struct{}
}

// NewScopeDB : ScopeDB 객체 할당 및 반환
func NewScopeDB(config *conf.Config, root *Repositories) (IRepository, error) {
	cfg := config.Repositories["contract-db"]
	_ = cfg

	// Build the data source name (dsn) string
	dsn := "user:password@tcp(127.0.0.1:3306)/databasename"

	// Open a new database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//r := &ContractDB{
	//	start: make(chan struct{}),
	//}
	//var err error
	//credential := options.Credential{
	//	Username: cfg["username"].(string),
	//	Password: cfg["pass"].(string),
	//}
	//
	//if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(cfg["datasource"].(string)).SetAuth(credential)); err != nil {
	//	return nil, err
	//} else if err := r.client.Ping(context.Background(), nil); err != nil {
	//	return nil, err
	//} else {
	//	db := r.client.Database(cfg["db"].(string))
	//	r.colContract = db.Collection("contracts")
	//	r.colExtContract = db.Collection("external_contracts")
	//	r.colEventManage = db.Collection("event_manage")
	//
	//}
	//
	//elog.Trace("load repository : ContractDB")
	//return r, nil
	return nil, nil
}
