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

// Repositories repositories manager
type Repositories struct {
	lock  sync.RWMutex
	conf  *conf.Config
	elems map[reflect.Type]reflect.Value
}

// RepositoryConstructor repository의 생성 함수 타입
type RepositoryConstructor func(conf *conf.Config, root *Repositories) (IRepository, error)

// NewRepositories 모든 repository를 생성 및 등록
func NewRepositories(cfg *conf.Config) (*Repositories, error) {
	r := &Repositories{
		conf:  cfg,
		elems: make(map[reflect.Type]reflect.Value),
	}

	for _, c := range []struct {
		constructor RepositoryConstructor
		config      *conf.Config
	}{
		{NewScopeDB, cfg},
	} {
		if err := r.Register(c.constructor, c.config); err != nil {
			return nil, err
		}
	}

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

// Register respository의 constructor를 호출하여 리턴된 instance를 map에 삽입한다.
func (p *Repositories) Register(constructor RepositoryConstructor, config *conf.Config) error {
	if r, err := constructor(config, p); err != nil {
		return err
	} else if r != nil {
		p.lock.Lock()
		defer p.lock.Unlock()

		if _, ok := p.elems[reflect.TypeOf(r)]; ok == true {
			return fmt.Errorf("duplicated instance of %v", reflect.TypeOf(r))
		} else {
			p.elems[reflect.TypeOf(r)] = reflect.ValueOf(r)
		}
	}
	return nil
}

func InitializeRepositories(cfg *conf.Config) (*Repositories, error) {
	return NewRepositories(cfg)
}

func (p *Repositories) GetScopeDB() *ScopeDB {
	p.lock.RLock()
	defer p.lock.RUnlock()

	v, ok := p.elems[reflect.TypeOf(&ScopeDB{})]
	if !ok {
		return nil
	}
	return v.Interface().(*ScopeDB)
}
