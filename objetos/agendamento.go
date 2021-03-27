package objetos

import "time"

type Agendamento struct {
	AID   uint64    `json:"aid"`
	ADATA time.Time `json:"data"`
}
