module bittoCralwer

go 1.19

require (
	github.com/ethereum/go-ethereum v1.13.1
	github.com/go-sql-driver/mysql v1.7.1
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

require (
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/holiman/uint256 v1.2.3 // indirect
	github.com/naoina/go-stringutil v0.1.0 // indirect
)

replace go-common/ether => ../go-common/ether
