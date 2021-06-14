package objetos

type Usuario struct {
	UID         uint64
	UNOME       string `json:"nome" validate:"required"`
	UEMAIL      string `json:"email" validate:"email,required"`
	UPASSWORD   string `json:"password" validate:"required"`
	UTOKEN      string
	UCPF        string `json:"cpf" validate:"cpf,required"`
	UNASCIMENTO int64  `json:"nascimento" validate:"idade,required"`
	UTELEFONE   uint64 `json:"telefone" validate:"required"`
	USEXO       string `json:"sexo" validate:"required"`
	CPID        uint64 `json:"cpid" validate:"required"`
}

type Dependente struct {
	DID         uint64 `json:"did"`
	DNOME       string `json:"nome" validate:"required"`
	DNASCIMENTO int64  `json:"nascimento" validate:"unix-passado,required"`
	DSEXO       string `json:"sexo" validate:"sexo,required"`
}
