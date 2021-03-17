package objetos

import "time"

type ItemEstoque struct {
	EID          uint32
	PID          uint32
	EQUANTIDADE  float32
	EVALORCOMPRA float32
	EVALORVENDA  float32
	EVALIDADE    time.Time
}
