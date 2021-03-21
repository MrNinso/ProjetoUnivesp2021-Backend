package objetos

type Medico struct {
	MID      uint64 `json:"mid"  validate:"required"`
	MNOME    string `json:"nome" validate:"required"`
	EID      uint   `json:"eid"  validate:"required"`
	HID      uint   `json:"hid"  validate:"required"`
	MATIVADO bool
}
