package objetos

import "time"

type Agendamento struct {
	AID   uint64    `json:"aid,omitempty"`
	ADATA time.Time `json:"data"`
}
