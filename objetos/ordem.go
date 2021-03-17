package objetos

import "time"

type Ordem struct {
	OID         uint32
	PID         uint32
	ODATA       time.Time
	OTIPO       string
	OQUANTIDADE float32
	OVALOR      float32
}
