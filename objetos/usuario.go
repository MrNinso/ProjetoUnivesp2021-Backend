package objetos

type Usuario struct {
	UID          uint64
	UNOME        string `json:"nome" validate:"alpha,required"`
	UEMAIL       string `json:"email" validate:"email,required"`
	UPASSWORD    string `json:"password" validate:"required"`
	UTOKEN       string
	UCPF         string `json:"cpf" validator:"numeric,cpf"`
	UUF          string `json:"uf" validator:"uf"`
	UCIDADE      string `json:"cidade" validator:"alpha,required"`
	UCEP         string `json:"cep" validate:"cep"`
	UENDERECO    string `json:"endereco" validate:"required"`
	UCOMPLEMENTO string `json:"complemento"`
	UATIVADO     bool   `json:"isAtivado"`
}
