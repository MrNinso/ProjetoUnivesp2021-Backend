package objetos

type Convenio struct {
	CID   uint   `json:"cid"`
	CNOME string `json:"nome"`
}

type Planos struct {
	CID    uint   `json:"cid"`
	CPID   uint64 `json:"cpid"`
	CPNOME string `json:"nome"`
}
