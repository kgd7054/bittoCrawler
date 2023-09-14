package model

import (
	conf "bittoCralwer/config"
	"fmt"
	"reflect"
	"sync"
)

// IRepository repository 타입
type IRepository interface {
	Start() error
}

// RepositoryConstructor repository의 생성 함수 타입
type RepositoryConstructor func(conf *conf.Config, root *Repositories) (IRepository, error)

// NewRepositories 모든 repository를 생성 및 등록
func NewRepositories(cfg *conf.Config) (*Repositories, error) {
	r := &Repositories{
		conf:  cfg,
		elems: make(map[reflect.Type]reflect.Value),
	}

	//for _, c := range []struct{
	//	constructor RepositoryConstructor
	//	config *conf.Config
	//}

	//for _, c := range []struct {
	//	constructor RepositoryConstructor
	//	config      *conf.Config
	//}{
	//	{NewRedisDB, cfg}, //다른 respository로서 제일먼저 추가되어야함.
	//	{NewAuthRedis, cfg},
	//} {
	//	if err := r.Register(c.constructor, c.config); err != nil {
	//		return nil, err
	//	}
	//}

	if err := func() error {
		r.lock.Lock()
		defer r.lock.Unlock()

		for t, e := range r.elems {
			if err := e.Interface().(IRepository).Start(); err != nil {
				//elog.Error("NewRepositories", "repository", t, "error", err)
				fmt.Println("NewRepositories", "repository", t, "error", err)
				return err
			}
		}
		return nil
	}(); err != nil {
		return nil, err
	}
	return r, nil
}

// Repositories repositories manager
type Repositories struct {
	lock  sync.RWMutex
	conf  *conf.Config
	elems map[reflect.Type]reflect.Value
}
