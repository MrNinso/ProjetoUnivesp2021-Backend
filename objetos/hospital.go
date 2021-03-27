package objetos

type Hospital struct {
	HID              int
	HNOME            string `json:"nome" validate:"required"`
	HUF              string `json:"uf" validate:"uf"`
	HCIDADE          string `json:"cidade" validate:"required"`
	HCEP             string `json:"cep" validate:"cep"`
	HENDERECO        string `json:"endereco" validate:"required"`
	HCOMPLEMENTO     string `json:"complemento" validate:"required"`
	HTELEFONE        uint64 `json:"telefone" validate:"required"`
	HISPRONTOSOCORRO bool   `json:"isProntoSocorro"`
	HATIVADO         bool
}