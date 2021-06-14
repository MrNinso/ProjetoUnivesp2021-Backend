package objetos

type Agendamento struct {
	AID   uint64 `json:"aid,omitempty"`
	DID   uint64 `json:"did,omitempty"`
	MID   uint64 `json:"mid,omitempty"`
	ADATA int64  `json:"data"`
}
