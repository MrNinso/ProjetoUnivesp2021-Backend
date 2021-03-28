package objetos

type Usuario struct {
	UID          uint64
	UNOME        string `json:"nome" validate:"required"`
	UEMAIL       string `json:"email" validate:"email,required"`
	UPASSWORD    string `json:"password" validate:"required"`
	UTOKEN       string
	UCPF         string `json:"cpf" validate:"cpf,required"`
	UUF          string `json:"uf" validate:"uf,required"`
	UCIDADE      string `json:"cidade" validate:"required"`
	UCEP         string `json:"cep" validate:"cep,required"`
	UENDERECO    string `json:"endereco" validate:"required"`
	UCOMPLEMENTO string `json:"complemento"`
	UATIVADO     bool
}
