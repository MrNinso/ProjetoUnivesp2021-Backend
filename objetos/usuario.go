package objetos

type Usuario struct {
	UID       int
	LID       int
	UNAME     string `json:"name" validate:"alpha,required"`
	UEMAIL    string `json:"email" validate:"email,required"`
	UPASSWORD string `json:"password" validate:"required"`
	UTOKEN    string
	UISADMIN  bool `json:"isAdmin"`
	ATIVADO   bool `json:"isAtivado"`
}
