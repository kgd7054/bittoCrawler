package dto

type EthereumBlock struct {
	BaseFeePerGas    string     `json:"baseFeePerGas"`
	Difficulty       string     `json:"difficulty"`
	ExtraData        string     `json:"extraData"`
	GasLimit         string     `json:"gasLimit"`
	GasUsed          string     `json:"gasUsed"`
	Hash             string     `json:"hash"`
	LogsBloom        string     `json:"logsBloom"`
	Miner            string     `json:"miner"`
	MixHash          string     `json:"mixHash"`
	Nonce            string     `json:"nonce"`
	Number           string     `json:"number"`
	ParentHash       string     `json:"parentHash"`
	ReceiptsRoot     string     `json:"receiptsRoot"`
	Sha3Uncles       string     `json:"sha3Uncles"`
	Size             string     `json:"size"`
	StateRoot        string     `json:"stateRoot"`
	Timestamp        string     `json:"timestamp"`
	TotalDifficulty  string     `json:"totalDifficulty"`
	Transactions     []string   `json:"transactions"`
	TransactionsRoot string     `json:"transactionsRoot"`
	Uncles           []string   `json:"uncles"`
	Withdrawal       Withdrawal `json:"withdrawal"`
	WithdrawalsRoot  string     `json:"withdrawalsRoot"`
}

type Withdrawal struct {
	Address        string `json:"address"`
	Amount         string `json:"amount"`
	Index          string `json:"index"`
	ValidatorIndex uint32 `json:"validatorIndex"`
}
