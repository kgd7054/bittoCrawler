package conf

import (
	"bittoCralwer/common"
	"github.com/naoina/toml"
	"os"
	"path"
)

//type Admin struct {
//	Address string
//	Pass    string
//	Desc    string
//}
//
//type Contract struct {
//	Name  string
//	Chain string
//	Owner string
//	Pass  string
//}

type Config struct {
	Datadir struct {
		Root     string
		Keystore string
		Journal  string
		Log      string
	}

	Repositories map[string]map[string]interface{}

	Port struct {
		Server string
		Http   int
	}

	Alchemy struct {
		ApiKey string
	}
}

func NewConfig(file string) *Config {
	c := new(Config)

	if file, err := os.Open(file); err != nil {
		panic(err)
	} else {
		defer file.Close()
		if err := toml.NewDecoder(file).Decode(c); err != nil {
			panic(err)
		} else {
			c.sanitize()
			return c
		}
	}
}

func (p *Config) sanitize() {
	if p.Datadir.Root[0] == byte('~') {
		p.Datadir.Root = path.Join(common.HomeDir(), p.Datadir.Root[1:])
	}
}
