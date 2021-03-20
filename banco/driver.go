package banco

import "github.com/MrNinso/ProjetoUnivesp2021-Backend/objetos"

type DriverBancoDados interface {
	Login(lid uint, uemail, upassword string) string

	IsValidToken(uemail, utoken string) (isValid, isAdmin bool)

	Logoff(lid uint, uemail, token string) uint

	CadastarUsuario(token string, usuario objetos.Usuario) uint

	AtualizarUsuario(token string, usuario objetos.Usuario) uint

	ListarUsuarios(token string, page uint8) []objetos.Usuario

	//TODO ADD CRIAR LOJA

	CriarProduto(token string, produto objetos.Produto) uint

	ListarProdutos(token string, page uint8) []objetos.Produto

	AtualizarProduto(token string, produto objetos.Produto) uint

	RegistrarOrdem(token string, ordem objetos.Ordem) uint

	ListarEstoque(token string, page uint8) []objetos.ItemEstoque
}
